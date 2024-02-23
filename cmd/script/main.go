package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit"
	"go.uber.org/zap"
)

const baseURL = "http://localhost:8080/api/accounts"

var logger *zap.Logger

func main() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	logger = l
	a := Accounts{}

	go func() {
		for {
			for i := 0; i < 3; i++ {
				go CreateAccount(&a)
			}
			time.Sleep(time.Second / 3)
		}
	}()

	time.Sleep(time.Second)
	ch := make(chan string, 7*8)
	go DepositAccount(ch)
	go func() {
		for {
			ch <- a.GetRandomID(3)
		}
	}()

	chTr := make(chan []string, 15*8)
	go Transfer(chTr)
	go func() {
		for {
			arr := []string{
				a.GetRandomID(5),
				a.GetRandomID(5),
			}

			if arr[0] == arr[1] {
				continue
			}

			chTr <- arr
		}
	}()

	for {
	}
}

func Transfer(ch chan []string) {
	for i := 0; i < 15*8; i++ {
		go func() {
			for arr := range ch {
				payer, payee := arr[0], arr[1]
				payload := map[string]any{
					"value": float64(rand.Intn(100000)) + rand.Float64(),
					"payee": payee,
				}
				res, _ := request(context.Background(), http.MethodPost, fmt.Sprintf("%s/%s/transfer", baseURL, payer), payload)
				var data struct {
					ID string `json:"transaction_id"`
				}
				res.Bind(&data)
				logger.Info("Transfer OK", zap.String("transaction_id", data.ID), zap.String("payer", payer), zap.String("payee", payee))
			}
		}()
	}
}

func DepositAccount(ch chan string) {
	for i := 0; i < 7*8; i++ {
		go func() {
			for id := range ch {
				rand.Seed(time.Now().UnixNano())
				payload := map[string]float64{
					"value": float64(rand.Intn(100000)) + rand.Float64(),
				}
				res, _ := request(context.Background(), http.MethodPost, fmt.Sprintf("%s/%s/deposit", baseURL, id), payload)

				var data struct {
					ID string `json:"transaction_id"`
				}
				res.Bind(&data)
				logger.Info("Deposit OK", zap.String("transaction_id", data.ID), zap.String("id", id))
			}
		}()
	}

}

func CreateAccount(a *Accounts) {
	payload := map[string]string{
		"customer_name":   gofakeit.Name(),
		"email":           gofakeit.Email(),
		"password":        gofakeit.Password(true, true, true, true, false, 20),
		"account_type":    "PERSONAL",
		"document_number": fmt.Sprintf("%d", gofakeit.Number(11111111111, 99999999999)),
		"phone_number":    gofakeit.Phone(),
	}

	res, err := request(context.Background(), http.MethodPost, baseURL, payload)
	if err != nil {
		logger.Error("error in create account", zap.Error(err))
		return
	}

	var data struct {
		ID string `json:"account_id"`
	}

	if err := res.Bind(&data); err != nil {
		logger.Error("error in create account", zap.Error(err))
		return
	}

	a.AppendID(data.ID)
	logger.Info("Account created", zap.String("account_id", data.ID), zap.String("name", payload["customer_name"]), zap.String("email", payload["email"]))
}

func request(ctx context.Context, method, url string, payload any) (*ClientResponse, error) {
	var body io.Reader
	if payload != nil {
		raw, _ := json.Marshal(payload)
		body = bytes.NewBuffer(json.RawMessage(raw))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return NewClientResponse(res)
}

type ClientResponse struct {
	body       []byte
	statusCode int
	headers    http.Header
}

func NewClientResponse(response *http.Response) (*ClientResponse, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return &ClientResponse{
		statusCode: response.StatusCode,
		headers:    response.Header,
		body:       body,
	}, nil
}

func (c *ClientResponse) Bind(v interface{}) error {
	return json.Unmarshal(c.body, v)
}

type Accounts struct {
	mu  sync.Mutex
	ids []string
}

func (a *Accounts) AppendID(id string) {
	a.mu.Lock()
	a.ids = append(a.ids, id)
	a.mu.Unlock()
}

func (a *Accounts) GetRandomID(k int) string {
	a.mu.Lock()
	defer a.mu.Unlock()
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(a.ids) / k)
	return a.ids[idx]
}
