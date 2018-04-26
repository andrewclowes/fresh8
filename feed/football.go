package feed

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Event represents an event from the feed
type Event struct {
	ID      *int       `json:"id,omitempty"`
	Name    *string    `json:"name,omitempty"`
	Time    *time.Time `json:"time,omitempty"`
	Markets *[]int     `json:"markets,omitempty"`
}

// Option represents a set of options for a market
type Option struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Odds *string `json:"odds,omitempty"`
}

// Market represents a market from the feed
type Market struct {
	ID      *string   `json:"id,omitempty"`
	Type    *string   `json:"type,omitempty"`
	Options *[]Option `json:"options,omitempty"`
}

// FootballService handles communication with the Football related data
type FootballService struct {
	client *Client
}

// ListEventIds lists the event IDs
func (s *FootballService) ListEventIds(ctx context.Context) (*[]int, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "football/events", nil)
	if err != nil {
		return nil, nil, err
	}
	i := new([]int)
	resp, err := s.client.Do(ctx, req, i)
	if err != nil {
		return nil, nil, err
	}
	return i, resp, nil
}

// GetEvent gets the event for a given id
func (s *FootballService) GetEvent(ctx context.Context, id int) (*Event, *http.Response, error) {
	u := fmt.Sprintf("football/events/%v", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	e := new(Event)
	resp, err := s.client.Do(ctx, req, e)
	if err != nil {
		return nil, nil, err
	}
	return e, resp, nil
}

// GetMarket gets the market for a given id
func (s *FootballService) GetMarket(ctx context.Context, id int) (*Market, *http.Response, error) {
	u := fmt.Sprintf("football/markets/%v", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	e := new(Market)
	resp, err := s.client.Do(ctx, req, e)
	if err != nil {
		return nil, nil, err
	}
	return e, resp, nil
}
