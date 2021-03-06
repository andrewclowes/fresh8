package feed

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/andrewclowes/fresh8/apiutil"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	client, _ = NewClient(server.URL, nil)

	return client, mux, server.URL, server.Close
}

func TestFootballService_ListEventIds(t *testing.T) {
	feedClient, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/football/events", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("Request method: %v, want %v", got, want)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[1,2,3]`))
	})

	got, _, err := feedClient.Football.ListEventIds(context.Background())
	if err != nil {
		t.Errorf("ListEventIds returned error: %v", err)
	}
	if want := []json.Number{"1", "2", "3"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ListEventIds = %+v, want %+v", got, want)
	}
}

func TestFootballService_GetEvent(t *testing.T) {
	feedClient, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/football/events/1", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("Request method: %v, want %v", got, want)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":1,"name":"A v B","time":"2018-04-25T12:00:00Z","markets":[1,2]}`))
	})

	got, _, err := feedClient.Football.GetEvent(context.Background(), 1)
	if err != nil {
		t.Errorf("GetEvent returned error: %v", err)
	}
	want := &Event{ID: apiutil.JSONNumber("1"), Name: apiutil.String("A v B"), Time: apiutil.String("2018-04-25T12:00:00Z"), Markets: &[]json.Number{"1", "2"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetEvent = %+v, want %+v", got, want)
	}
}

func TestFootballService_GetMarket_SingleOption(t *testing.T) {
	feedClient, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/football/markets/1", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("Request method: %v, want %v", got, want)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"1","type":"win","options":{"id":"1","name":"win","odds":"1/2"}}`))
	})

	got, _, err := feedClient.Football.GetMarket(context.Background(), 1)
	if err != nil {
		t.Errorf("GetMarket returned error: %v", err)
	}
	want := &Market{ID: apiutil.JSONNumber("1"), Type: apiutil.String("win"), Options: &Options{Option{ID: apiutil.JSONNumber("1"), Name: apiutil.String("win"), Odds: apiutil.String("1/2")}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetMarket = %+v, want %+v", got, want)
	}
}

func TestFootballService_GetMarket_MultipleOptions(t *testing.T) {
	feedClient, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/football/markets/1", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("Request method: %v, want %v", got, want)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"1","type":"win","options":[{"id":"1","name":"win","odds":"1/2"}]}`))
	})

	got, _, err := feedClient.Football.GetMarket(context.Background(), 1)
	if err != nil {
		t.Errorf("GetMarket returned error: %v", err)
	}
	want := &Market{ID: apiutil.JSONNumber("1"), Type: apiutil.String("win"), Options: &Options{Option{ID: apiutil.JSONNumber("1"), Name: apiutil.String("win"), Odds: apiutil.String("1/2")}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetMarket = %+v, want %+v", got, want)
	}
}
