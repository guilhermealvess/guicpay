package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/guilhermealvess/guicpay/domain/entity"
	"github.com/guilhermealvess/guicpay/domain/gateway"
	clienthttp "github.com/guilhermealvess/guicpay/internal/client_http"
	"go.opentelemetry.io/otel"
)

type authorizationService struct {
	client clienthttp.HTTPClient
}

func NewAuthorizationService(baseURL string) gateway.AuthorizationService {
	return &authorizationService{
		client: clienthttp.NewHTTPClient(baseURL),
	}
}

func (s *authorizationService) Authorize(ctx context.Context, account entity.Account) error {
	ctx, span := otel.GetTracerProvider().Tracer("my-server").Start(ctx, "AuthorizationService.Auth")
	defer span.End()

	const endpoint = "/auth"
	res, err := s.client.Request(ctx, http.MethodPost, endpoint, clienthttp.WithPayload(account))
	if err != nil {
		return err
	}

	if err := res.Error(); err != nil {
		return err
	}

	var data map[string]interface{}
	if err := res.Bind(&data); err != nil {
		return fmt.Errorf("TODO: ... %w", err)
	}

	if data["message"] != true {
		return errors.New("TODO:")
	}

	return nil
}
