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
	"go.opentelemetry.io/otel"
)

type authorizationService struct {
	client  clienthttp.ClientHttp
	baseURL string
}

func NewAuthorizationService(baseURL string) gateway.AuthorizationService {
	return &authorizationService{
		client:  clienthttp.NewClient(),
		baseURL: baseURL,
	}
}

func (s *authorizationService) Authorize(ctx context.Context, account entity.Account) error {
	ctx, span := otel.GetTracerProvider().Tracer("my-server").Start(ctx, "Authorize")
	defer span.End()

	url := s.baseURL + "/auth"
	res, err := s.client.Send(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("TODO: ... %w", err)
	}

	if data["message"] != true {
		return errors.New("TODO:")
	}

	return nil
}
