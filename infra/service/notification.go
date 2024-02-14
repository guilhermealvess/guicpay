package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/gateway"
	clienthttp "github.com/guilhermealvess/guicpay/internal/client_http"
)

type notificationService struct {
	clientHttp clienthttp.ClientHttp
	baseURL    string
}

func NewNotificationService(baseURL string) gateway.NotificationService {
	return &notificationService{
		clientHttp: clienthttp.NewClient(),
		baseURL:    baseURL,
	}
}

func (s *notificationService) Notify(ctx context.Context, account entity.Account, transaction entity.Transaction) error {
	url := s.baseURL + "/dispatch"
	payload := map[string]string{
		"message": fmt.Sprintf("%s, você recebeu uma nova transferência no valor de %s", account.CustomerName, transaction.Amount.String()),
	}
	raw, _ := json.Marshal(payload)
	res, err := s.clientHttp.Send(ctx, http.MethodGet, url, nil, raw)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("TODO: ... %w", err)
	}

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("TODO: ..., %w", err)
	}

	if data["message"].(string) != "TODO:" {
		return errors.New("TODO:")
	}

	return nil
}
