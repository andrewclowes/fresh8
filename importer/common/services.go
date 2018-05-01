package common

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/andrewclowes/fresh8/feed"
	"github.com/andrewclowes/fresh8/store"
)

// FootballService defines the football service for the feed api
type FootballService interface {
	ListEventIds(ctx context.Context) ([]json.Number, *http.Response, error)
	GetEvent(ctx context.Context, id int) (*feed.Event, *http.Response, error)
	GetMarket(ctx context.Context, id int) (*feed.Market, *http.Response, error)
}

// StoreService defines the event service for the store api
type StoreService interface {
	Create(ctx context.Context, event *store.Event) (*http.Response, error)
}
