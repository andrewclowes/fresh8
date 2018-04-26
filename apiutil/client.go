package apiutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ClientBase manages communications via http
type ClientBase struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

// NewClientBase creates a new instance of an http channel
func NewClientBase(baseURL string, httpClient *http.Client) (*ClientBase, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &ClientBase{httpClient: httpClient}
	c.BaseURL = u
	return c, nil
}

// NewRequest creates an API request
func (b *ClientBase) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := b.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if b.UserAgent != "" {
		req.Header.Set("User-Agent", b.UserAgent)
	}
	return req, nil
}

// Do sends an API request and returns the API response
func (b *ClientBase) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	r := req.WithContext(ctx)
	resp, err := b.httpClient.Do(r)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()
	err = checkResponse(resp)
	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

// ErrorResponse reports errors caused by an API response
type ErrorResponse struct {
	Response *http.Response
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode)
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; c < 200 || c >= 400 {
		return &ErrorResponse{Response: r}
	}
	return nil
}
