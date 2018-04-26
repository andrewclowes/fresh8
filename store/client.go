package store

import (
	"net/http"

	"github.com/andrewclowes/fresh8/apiutil"
)

// Client manages communication with the store API
type Client struct {
	*apiutil.ClientBase

	Event *EventService
}

// NewClient creates a new instance of the client for the store api
func NewClient(baseURL string, httpClient *http.Client) (*Client, error) {
	b, err := apiutil.NewClientBase(baseURL, httpClient)
	if err != nil {
		return nil, err
	}
	c := &Client{ClientBase: b}
	c.Event = &EventService{client: c}
	return c, nil
}
