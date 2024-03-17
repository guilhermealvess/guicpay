package clienthttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type HTTPClient interface {
	Request(ctx context.Context, method string, url string, options ...RequestOptions) (*Response, error)
}

type RequestOptions func(req *http.Request)

func WithPayload(v any) RequestOptions {
	return func(req *http.Request) {
		body, _ := json.Marshal(v)
		req.Body = io.NopCloser(bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
}

type httpClient struct {
	baseURL string
}

func NewHTTPClient(baseURL string) HTTPClient {
	return &httpClient{
		baseURL: baseURL,
	}
}

func (c *httpClient) Request(ctx context.Context, method string, url string, options ...RequestOptions) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+url, nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range options {
		opt(req)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Content:  body,
		Response: res,
	}, nil
}

type Response struct {
	Content  []byte
	Response *http.Response
}

func (r *Response) Bind(v interface{}) error {
	return json.Unmarshal(r.Content, v)
}

func (r *Response) Error() error {
	if r.Response.StatusCode >= 400 && r.Response.StatusCode < 500 {
		return errors.New("client error")
	}

	if r.Response.StatusCode >= 500 {
		return errors.New("server error")
	}

	return nil
}
