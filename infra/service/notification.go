package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/gateway"
	clienthttp "github.com/guilhermealvess/guicpay/internal/client_http"
	"github.com/guilhermealvess/guicpay/internal/logger"
	"go.opentelemetry.io/otel"
)

var notificationMockServer *httptest.Server

type notificationService struct {
	clientHttp clienthttp.HTTPClient
}

func NewNotificationService(baseURL string) gateway.NotificationService {
	c := clienthttp.NewHTTPClient(baseURL)
	if os.Getenv("USE_MOCK_SERVER") == "true" {
		c = clienthttp.NewHTTPClient(notificationMockServer.URL)
	}

	return &notificationService{
		clientHttp: c,
	}
}

func (s *notificationService) Notify(ctx context.Context, account entity.Account, transaction entity.Transaction) error {
	ctx, span := otel.GetTracerProvider().Tracer("my-server").Start(ctx, "NotificationService.Notify")
	defer span.End()

	const endpoint = "/dispatch"
	payload := map[string]string{
		"message": fmt.Sprintf("%s, você recebeu uma nova transferência no valor de %s", account.CustomerName, transaction.Amount.String()),
	}

	res, err := s.clientHttp.Request(ctx, http.MethodPost, endpoint, clienthttp.WithPayload(payload))
	if err != nil {
		span.RecordError(err)
		return err
	}

	if err := res.Error(); err != nil {
		span.RecordError(err)
		return err
	}

	var data struct {
		Message string `json:"message"`
	}

	if err := res.Bind(&data); err != nil {
		span.RecordError(err)
		return fmt.Errorf("notification error: %w", err)
	}

	if data.Message != "Autorizado" {
		span.RecordError(errors.New(data.Message))
		return errors.New(data.Message)
	}

	return nil
}

func init() {
	notificationMockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(json.RawMessage(`{"message": "Autorizado"}`))

		logger.Logger.Info("Notification mock server")
	}))
}
