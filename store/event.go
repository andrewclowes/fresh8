package store

import (
	"context"
	"net/http"
)

// Event represents an event from the store
type Event struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Time    string   `json:"time"`
	Markets []Market `json:"markets"`
}

// Market represents a market from the store
type Market struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Options []Option `json:"options"`
}

// Option represents a set of options for a market
type Option struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Num  int    `json:"num"`
	Den  int    `json:"den"`
}

// EventService handles communication with the Event related data
type EventService struct {
	client *Client
}

// Create creates an event in the store
func (s *EventService) Create(ctx context.Context, event *Event) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "event", event)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
