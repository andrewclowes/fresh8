package feed

import (
	"net/http"

	"github.com/andrewclowes/fresh8/apiutil"
)

// Client manages communication with the feed API
type Client struct {
	*apiutil.ClientBase

	Football *FootballService
}

// NewClient creates a new instance of the client for the feed api
func NewClient(baseURL string, httpClient *http.Client) (*Client, error) {
	b, err := apiutil.NewClientBase(baseURL, httpClient)
	if err != nil {
		return nil, err
	}
	c := &Client{ClientBase: b}
	c.Football = &FootballService{client: c}
	return c, nil
}
