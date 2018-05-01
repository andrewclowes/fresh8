package common

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/andrewclowes/fresh8/feed"
	"github.com/andrewclowes/fresh8/store"
)

// FootballEventIdsService defines the football event ids provider for the feed api
type FootballEventIdsService interface {
	ListEventIds(ctx context.Context) ([]json.Number, *http.Response, error)
}

// FootballEventService defines the football event ids providerfor the feed api
type FootballEventService interface {
	GetEvent(ctx context.Context, id int) (*feed.Event, *http.Response, error)
}

// FootballMarketService defines the football event ids service for the feed api
type FootballMarketService interface {
	GetMarket(ctx context.Context, id int) (*feed.Market, *http.Response, error)
}

// FootballService defines the football service for the feed api
type FootballService interface {
	FootballEventIdsService
	FootballEventService
	FootballMarketService
}

// StoreEventService defines the store event service for the store api
type StoreEventService interface {
	Create(ctx context.Context, event *store.Event) (*http.Response, error)
}

// StoreService defines the store service for the store api
type StoreService interface {
	StoreEventService
}
