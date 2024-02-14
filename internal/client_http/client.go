package clienthttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type ClientHttp interface {
	Send(ctx context.Context, method, url string, query url.Values, bodyJSON json.RawMessage, pathParams ...string) (*http.Response, error)
}

type clientHttp struct {
	client http.Client
}

func NewClient() ClientHttp {
	return &clientHttp{
		client: *http.DefaultClient,
	}
}

func (c *clientHttp) Send(ctx context.Context, method, url string, query url.Values, bodyJSON json.RawMessage, pathParams ...string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return c.client.Do(req)
}

func (r *clientHttp) Bind(res *http.Response, v interface{}) error {
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}
