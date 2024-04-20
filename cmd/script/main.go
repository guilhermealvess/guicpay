package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/guilhermealvess/guicpay/domain/usecase"
	clienthttp "github.com/guilhermealvess/guicpay/internal/client_http"
	"go.uber.org/zap"
)

type processor struct {
	client clienthttp.HTTPClient
}

var (
	ids    = make([]Account, 0)
	logger *zap.Logger
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	baseURL := os.Getenv("API_URL")
	if baseURL == "" {
		baseURL = "https://api.guicpay.tech"
	}

	p := &processor{
		client: clienthttp.NewHTTPClient(baseURL),
	}

	logger, _ = zap.NewProduction()
	cicle := make(Cicle, 0)

	for range 10 {
		account, err := p.CreateAccount()
		if err != nil {
			logger.Error("Error creating account", zap.Error(err))
			return
		}

		ids = append(ids, *account)

		go p.Deposit(*account)

		cicle = append(cicle, group{Account: *account})
	}

	cicle.Mount()

	for _, it := range cicle {
		go p.Transfer(it.Account, *it.Next)
		go p.Transfer(it.Account, *it.Previous)
	}

	time.Sleep(time.Minute * 2)
}

type Account struct {
	ID    string `json:"account_id"`
	Token string `json:"token"`
}

type group struct {
	Account  Account
	Next     *Account
	Previous *Account
}

type Cicle []group

func (c Cicle) Mount() {
	for i := range len(c) {
		if i == 0 {
			c[i].Next = &c[i+1].Account
			c[i].Previous = &c[len(c)-1].Account
			continue
		}

		if i == len(c)-1 {
			c[i].Next = &c[0].Account
			c[i].Previous = &c[i-1].Account
			continue
		}

		c[i].Next = &c[i+1].Account
		c[i].Previous = &c[i-1].Account
	}
}

func (p *processor) CreateAccount() (*Account, error) {
	payload := usecase.NewAccountInput{
		Name:           gofakeit.Name(),
		Email:          gofakeit.Email(),
		Password:       gofakeit.Password(true, true, true, true, false, 10),
		Type:           "PERSONAL",
		DocumentNumber: fmt.Sprintf("%d", gofakeit.Number(11111111111, 99999999999)),
		PhoneNumber:    gofakeit.Phone(),
	}

	res, err := p.client.Request(context.Background(), http.MethodPost, "/accounts", clienthttp.WithPayload(payload), clienthttp.WithUserAgent("guicpay-script"))
	if err != nil {
		return nil, err
	}

	if err := res.Error(); err != nil {
		return nil, fmt.Errorf(`[%w] body: "%s"`, err, res.Content)
	}

	requestID := res.Response.Header.Get("X-Request-ID")
	logger.Info("Create account", zap.String("request_id", requestID), zap.String("email", payload.Email), zap.Int("status_code", res.Response.StatusCode))

	if err := res.Error(); err != nil {
		return nil, err
	}

	var data struct {
		ID string `json:"account_id"`
	}

	if err := res.Bind(&data); err != nil {
		return nil, err
	}

	token, err := p.Authenticate(payload)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:    data.ID,
		Token: token,
	}, nil
}

func (p *processor) Authenticate(payload usecase.NewAccountInput) (string, error) {
	res, err := p.client.Request(context.Background(), http.MethodPost, "/auth", clienthttp.WithPayload(map[string]interface{}{"email": payload.Email, "password": payload.Password}), clienthttp.WithUserAgent("guicpay-script"))
	if err != nil {
		return "", err
	}

	requestID := res.Response.Header.Get("X-Request-ID")
	logger.Info("Authenticate account", zap.String("request_id", requestID), zap.String("email", payload.Email), zap.Int("status_code", res.Response.StatusCode))

	if err := res.Error(); err != nil {
		return "", err
	}

	var data struct {
		Token string `json:"token"`
	}

	if err := res.Bind(&data); err != nil {
		return "", err
	}

	return data.Token, nil
}

func (p *processor) Deposit(account Account) {
	for {
		payload := map[string]interface{}{
			"value": gofakeit.Price(1, 1000),
		}

		res, err := p.client.Request(context.Background(), http.MethodPost, "/transactions/deposit", clienthttp.WithPayload(payload), clienthttp.WithUserAgent("guicpay-script"), clienthttp.WithToken(account.Token))
		if err != nil {
			logger.Error("Error depositing account", zap.Error(err))
			return
		}

		requestID := res.Response.Header.Get("X-Request-ID")
		logger.Info("Deposit account", zap.String("request_id", requestID), zap.String("account_id", account.ID), zap.Int("status_code", res.Response.StatusCode))

		if err := res.Error(); err != nil {
			logger.Error("Error depositing account", zap.Error(err))
			return
		}

		var data struct {
			ID string `json:"transaction_id"`
		}

		if err := res.Bind(&data); err != nil {
			logger.Error("Error depositing account", zap.Error(err))
			return
		}

		logger.Info("Deposit account", zap.String("account_id", account.ID), zap.Int("status_code", res.Response.StatusCode))
	}
}

func (p *processor) Transfer(payer, payee Account) {
	for {
		payload := map[string]interface{}{
			"value": gofakeit.Price(1, 1000),
			"payee": payee.ID,
		}

		res, err := p.client.Request(context.Background(), http.MethodPost, "/transactions/transfer", clienthttp.WithPayload(payload), clienthttp.WithUserAgent("guicpay-script"), clienthttp.WithToken(payer.Token))
		if err != nil {
			logger.Error("Error transfering account", zap.Error(err))
		}

		requestID := res.Response.Header.Get("X-Request-ID")
		logger.Info("Transfer account", zap.String("request_id", requestID), zap.String("payer_id", payer.ID), zap.String("payee_id", payee.ID), zap.Int("status_code", res.Response.StatusCode))

		if err := res.Error(); err != nil {
			logger.Error("Error transfering account", zap.Error(err))
		}

		var data struct {
			ID string `json:"transaction_id"`
		}

		if err := res.Bind(&data); err != nil {
			logger.Error("Error transfering account", zap.Error(err))
		}

		logger.Info("Transfer account", zap.String("payer_id", payer.ID), zap.String("payee_id", payee.ID), zap.Int("status_code", res.Response.StatusCode))
	}
}
