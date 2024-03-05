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
		return fmt.Errorf("authorization_service: serialize data, %w", err)
	}

	if data["message"] != true {
		return errors.Join(errors.New("authorization_service: not authorized"), entity.ErrUnprocessableEntity)
	}

	return nil
}
