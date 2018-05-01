package store

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	client, _ = NewClient(server.URL, nil)

	return client, mux, server.URL, server.Close
}

func TestEventService_Create(t *testing.T) {
	storeClient, mux, _, teardown := setup()
	defer teardown()

	input := &Event{ID: "1", Name: "Test", Time: time.Date(2018, 4, 25, 12, 0, 0, 0, time.UTC), Markets: []Market{Market{ID: "1", Type: "Win", Options: []Option{Option{ID: "1", Name: "Win", Num: 1, Den: 2}}}}}

	mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		v := new(Event)
		json.NewDecoder(r.Body).Decode(v)

		if got, want := r.Method, "POST"; got != want {
			t.Errorf("Request method: %v, want %v", got, want)
		}

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusOK)
	})

	_, err := storeClient.Event.Create(context.Background(), input)
	if err != nil {
		t.Errorf("Create returned error: %v", err)
	}
}
