package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/guilhermealvess/guicpay/domain/usecase"
	clienthttp "github.com/guilhermealvess/guicpay/internal/client_http"
	"go.uber.org/zap"
)

var (
	client     clienthttp.HTTPClient
	accountMap sync.Map
	m          = &manager{}
	logger     *zap.Logger
)

const (
	baseURL = "http://localhost:8080/api/accounts"
	k       = 3
	n       = 50
)

func init() {
	client = clienthttp.NewHTTPClient("http://localhost:8080/api")
	accountMap = sync.Map{}
	logger, _ = zap.NewProduction()
}

type pair struct {
	payer, payee string
}

type manager struct {
	mu  sync.Mutex
	ids []string
}

func (m *manager) Register(id string) {
	m.mu.Lock()
	m.ids = append(m.ids, id)
	m.mu.Unlock()
}

func (m *manager) GetRandomID(v int) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	idx := randRange(0, int(len(m.ids)/v)+1)
	return m.ids[idx]
}

func main() {
	go func() {
		for {
			time.Sleep(time.Second / k)

			for range k {
				payload := usecase.NewAccountInput{
					Name:           gofakeit.Name(),
					Email:          gofakeit.Email(),
					Password:       gofakeit.Password(true, true, true, true, false, 10),
					Type:           "PERSONAL",
					DocumentNumber: fmt.Sprintf("%d", gofakeit.Number(11111111111, 99999999999)),
					PhoneNumber:    gofakeit.Phone(),
				}

				id, err := NewAccount(payload)
				if err != nil {
					logger.Error("Error creating account", zap.Error(err))
					continue
				}

				token, err := Authenticate(payload)
				if err != nil {
					logger.Error("Error authenticating account", zap.Error(err))
					continue
				}

				accountMap.Store(id, token)
				m.Register(id)
			}
		}
	}()

	time.Sleep(time.Second * 5)

	ch := make(chan string, n)
	go ExecuteDepositLoop(ch)
	go func() {
		for {
			ch <- m.GetRandomID(k)
		}
	}()

	chTransfer := make(chan pair, n)
	go ExecuteTransferLoop(chTransfer)
	go func() {
		for {
			arr := pair{payer: m.GetRandomID(k), payee: m.GetRandomID(k)}
			chTransfer <- arr
		}
	}()

	<-time.After(time.Minute)
}

func NewAccount(payload usecase.NewAccountInput) (string, error) {
	res, err := client.Request(context.Background(), http.MethodPost, "/accounts", clienthttp.WithPayload(payload))
	if err != nil {
		return "", err
	}

	if err := res.Error(); err != nil {
		return "", err
	}

	var data struct {
		ID string `json:"account_id"`
	}

	if err := res.Bind(&data); err != nil {
		return "", err
	}

	return data.ID, nil
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func ExecuteDepositLoop(ch <-chan string) {
	for range n {
		go func() {
			for id := range ch {
				payload := map[string]interface{}{
					"value": gofakeit.Number(100, 1000),
				}

				token, ok := accountMap.Load(id)
				if !ok {
					continue
				}

				res, err := client.Request(context.Background(), http.MethodPost, "/transactions/deposit", clienthttp.WithPayload(payload), clienthttp.WithToken(token.(string)))
				if err != nil {
					continue
				}

				if err := res.Error(); err != nil {
					continue
				}

				var data struct {
					ID string `json:"transaction_id"`
				}

				if err := res.Bind(&data); err != nil {
					continue
				}

				logger.Info("Deposit OK", zap.String("transaction_id", data.ID), zap.String("account_id", id))
			}
		}()
	}
}

func ExecuteTransferLoop(ch <-chan pair) {
	for range n {
		go func() {
			for pair := range ch {
				payload := map[string]interface{}{
					"value": gofakeit.Number(100, 1000),
					"payee": pair.payee,
				}

				token, ok := accountMap.Load(pair.payer)
				if !ok {
					continue
				}

				res, err := client.Request(context.Background(), http.MethodPost, "/transactions/transfer", clienthttp.WithPayload(payload), clienthttp.WithToken(token.(string)))
				if err != nil {
					continue
				}

				if err := res.Error(); err != nil {
					continue
				}

				var data struct {
					ID string `json:"transaction_id"`
				}

				if err := res.Bind(&data); err != nil {
					continue
				}

				logger.Info("Transfer OK", zap.String("transaction_id", data.ID), zap.String("payer", pair.payer), zap.String("payee", pair.payee))
			}
		}()
	}
}

func Authenticate(data usecase.NewAccountInput) (string, error) {
	payload := map[string]interface{}{
		"email":    data.Email,
		"password": data.Password,
	}

	res, err := client.Request(context.Background(), http.MethodPost, "/auth", clienthttp.WithPayload(payload))
	if err != nil {
		return "", err
	}

	if err := res.Error(); err != nil {
		return "", err
	}

	var token struct {
		Token string `json:"token"`
	}
	if err := res.Bind(&token); err != nil {
		return "", err
	}

	return token.Token, nil
}
