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
	"go.opentelemetry.io/otel"
)

var authorizationMockService *httptest.Server

type authorizationService struct {
	client clienthttp.HTTPClient
}

func NewAuthorizationService(baseURL string) gateway.AuthorizationService {
	c := clienthttp.NewHTTPClient(baseURL)
	if os.Getenv("USE_MOCK_SERVER") == "true" {
		c = clienthttp.NewHTTPClient(authorizationMockService.URL)
	}

	return &authorizationService{
		client: c,
	}
}

func (s *authorizationService) Authorize(ctx context.Context, account entity.Account) error {
	ctx, span := otel.GetTracerProvider().Tracer("my-server").Start(ctx, "AuthorizationService.Auth")
	defer span.End()

	const endpoint = "/auth"
	res, err := s.client.Request(ctx, http.MethodPost, endpoint, clienthttp.WithPayload(account))
	if err != nil {
		span.RecordError(err)
		return err
	}

	if err := res.Error(); err != nil {
		span.RecordError(err)
		return err
	}

	var data map[string]interface{}
	if err := res.Bind(&data); err != nil {
		span.RecordError(err)
		return fmt.Errorf("TODO: ... %w", err)
	}

	if data["message"] != true {
		err := errors.New("TODO:")
		span.RecordError(err)
		return err
	}

	return nil
}

func init() {
	authorizationMockService = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(json.RawMessage(`{"message": true}`))
	}))
}
