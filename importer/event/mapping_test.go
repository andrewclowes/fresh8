package event

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/andrewclowes/fresh8/apiutil"
	"github.com/andrewclowes/fresh8/feed"
	"github.com/andrewclowes/fresh8/store"
)

func TestEventMapper(t *testing.T) {
	input := &feed.Event{
		ID:      apiutil.JSONNumber("1"),
		Name:    apiutil.String("A v B"),
		Time:    apiutil.String("2018-04-25T12:00:00Z"),
		Markets: &[]json.Number{"2"},
	}

	mockTimeParser := func(raw string) (*time.Time, error) {
		t := time.Date(2018, 4, 25, 12, 0, 0, 0, time.UTC)
		return &t, nil
	}

	mapEvent := newEventMapper(mockTimeParser)

	got, err := mapEvent(input)

	if err != nil {
		t.Fatalf("mapEvent returned error: %v", err)
	}

	want := &store.Event{ID: "1", Name: "A v B", Time: time.Date(2018, 4, 25, 12, 0, 0, 0, time.UTC), Markets: []store.Market{store.Market{ID: "2"}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapEvent = %+v, want %+v", *got, want)
	}
}

func TestMarketMapper(t *testing.T) {
	input := &feed.Market{
		ID:   apiutil.JSONNumber("1"),
		Type: apiutil.String("win"),
		Options: &feed.Options{
			feed.Option{
				ID:   apiutil.JSONNumber("1"),
				Name: apiutil.String("win"),
				Odds: apiutil.String("1/2"),
			},
		},
	}

	mockOptionMapper := func(option *feed.Option) (*store.Option, error) {
		return &store.Option{ID: "1", Name: "win", Num: 1, Den: 2}, nil
	}

	mapMarket := newMarketMapper(mockOptionMapper)

	got, err := mapMarket(input)

	if err != nil {
		t.Fatalf("mapMarket returned error: %v", err)
	}

	want := &store.Market{ID: "1", Type: "win", Options: []store.Option{store.Option{ID: "1", Name: "win", Num: 1, Den: 2}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapMarket = %+v, want %+v", *got, want)
	}
}

func TestOptionMapper(t *testing.T) {
	input := &feed.Option{
		ID:   apiutil.JSONNumber("1"),
		Name: apiutil.String("win"),
		Odds: apiutil.String("1/2"),
	}

	mockOddsParser := func(raw string) (*odds, error) {
		return &odds{Num: 1, Den: 2}, nil
	}

	mapOption := newOptionMapper(mockOddsParser)

	got, err := mapOption(input)

	if err != nil {
		t.Fatalf("mapOption returned error: %v", err)
	}

	want := &store.Option{ID: "1", Name: "win", Num: 1, Den: 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapOption = %+v, want %+v", *got, want)
	}
}
