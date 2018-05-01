package event

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/andrewclowes/fresh8/apiutil"
	"github.com/andrewclowes/fresh8/feed"
	"github.com/andrewclowes/fresh8/store"
)

type MockFootballEventIds struct {
	MockListEventIds func(ctx context.Context) ([]json.Number, *http.Response, error)
}

func (m *MockFootballEventIds) ListEventIds(ctx context.Context) ([]json.Number, *http.Response, error) {
	if m.MockListEventIds != nil {
		return m.MockListEventIds(ctx)
	}
	return []json.Number{}, nil, nil
}

func TestGetEventIdsStep(t *testing.T) {
	want := []json.Number{"1", "2", "3"}

	mockListEventIds := &MockFootballEventIds{
		MockListEventIds: func(ctx context.Context) ([]json.Number, *http.Response, error) {
			return want, nil, nil
		},
	}

	runStep := newGetEventIdsStep(mockListEventIds)

	out := make(chan interface{})
	go func() {
		runStep(nil, out, nil)
		close(out)
	}()

	var got []json.Number
	for i := range out {
		n, _ := i.(json.Number)
		got = append(got, n)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("getEventIdsStep = %+v, want %+v", got, want)
	}
}

type MockFootballEvent struct {
	MockGetEvent func(ctx context.Context, id int) (*feed.Event, *http.Response, error)
}

func (m *MockFootballEvent) GetEvent(ctx context.Context, id int) (*feed.Event, *http.Response, error) {
	if m.MockGetEvent != nil {
		return m.MockGetEvent(ctx, id)
	}
	return &feed.Event{}, nil, nil
}

func TestGetEventStep(t *testing.T) {
	mockGetEvent := &MockFootballEvent{
		MockGetEvent: func(ctx context.Context, id int) (*feed.Event, *http.Response, error) {
			e := &feed.Event{ID: apiutil.JSONNumber("1"), Name: apiutil.String("A v B"), Time: apiutil.String("2018-04-25T12:00:00Z"), Markets: &[]json.Number{"2"}}
			return e, nil, nil
		},
	}

	want := &store.Event{ID: "1", Name: "A v B", Time: time.Date(2018, 4, 25, 12, 0, 0, 0, time.UTC), Markets: []store.Market{store.Market{ID: "2"}}}
	mockEventMapper := func(option *feed.Event) (*store.Event, error) {
		return want, nil
	}

	runStepHandler := newGetEventStep(mockGetEvent, mockEventMapper)
	got := runStepHandler(json.Number("1"), nil)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("getEventStep = %+v, want %+v", got, want)
	}
}

type MockFootballMarket struct {
	MockGetMarket func(ctx context.Context, id int) (*feed.Market, *http.Response, error)
}

func (m *MockFootballMarket) GetMarket(ctx context.Context, id int) (*feed.Market, *http.Response, error) {
	if m.MockGetMarket != nil {
		return m.MockGetMarket(ctx, id)
	}
	return &feed.Market{}, nil, nil
}

func TestGetMarketsStep(t *testing.T) {
	input := &store.Event{ID: "1", Name: "A v B", Time: time.Date(2018, 4, 25, 12, 0, 0, 0, time.UTC), Markets: []store.Market{store.Market{ID: "1"}}}

	mockGetMarket := &MockFootballMarket{
		MockGetMarket: func(ctx context.Context, id int) (*feed.Market, *http.Response, error) {
			m := &feed.Market{ID: apiutil.JSONNumber("1"), Type: apiutil.String("win"), Options: &feed.Options{feed.Option{ID: apiutil.JSONNumber("1"), Name: apiutil.String("win"), Odds: apiutil.String("1/2")}}}
			return m, nil, nil
		},
	}

	mockMarketMapper := func(option *feed.Market) (*store.Market, error) {
		m := &store.Market{ID: "1", Type: "win", Options: []store.Option{store.Option{ID: "1", Name: "win", Num: 1, Den: 2}}}
		return m, nil
	}

	runStepHandler := newGetMarketsStep(mockGetMarket, mockMarketMapper)
	got := runStepHandler(input, nil)

	want := &store.Event{ID: "1", Name: "A v B", Time: time.Date(2018, 4, 25, 12, 0, 0, 0, time.UTC), Markets: []store.Market{store.Market{ID: "1", Type: "win", Options: []store.Option{store.Option{ID: "1", Name: "win", Num: 1, Den: 2}}}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("getMarketsStep = %+v, want %+v", got, want)
	}
}

type MockStoreEventCreate struct {
	MockEventCreate func(ctx context.Context, event *store.Event) (*http.Response, error)
}

func (m *MockStoreEventCreate) Create(ctx context.Context, event *store.Event) (*http.Response, error) {
	if m.MockEventCreate != nil {
		return m.MockEventCreate(ctx, event)
	}
	return nil, nil
}

func TestStoreSendStep(t *testing.T) {
	want := &store.Event{ID: "1", Name: "A v B", Time: time.Date(2018, 4, 25, 12, 0, 0, 0, time.UTC), Markets: []store.Market{store.Market{ID: "1", Type: "win", Options: []store.Option{store.Option{ID: "1", Name: "win", Num: 1, Den: 2}}}}}

	var got *store.Event
	mockStoreSend := &MockStoreEventCreate{
		MockEventCreate: func(ctx context.Context, event *store.Event) (*http.Response, error) {
			got = event
			return nil, nil
		},
	}

	runStepHandler := newStoreSendStep(mockStoreSend)
	runStepHandler(want, nil)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("getMarketsStep = %+v, want %+v", got, want)
	}
}
