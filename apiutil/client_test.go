package apiutil

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type foo struct {
	A string
}

func setup() (client *ClientBase, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	client, _ = NewClientBase(server.URL, nil)

	return client, mux, server.URL, server.Close
}

func TestNewRequest(t *testing.T) {
	baseURL := "http://localhost/"
	c, _ := NewClientBase(baseURL, nil)
	c.UserAgent = "test"

	inURL, outURL := "/foo", baseURL+"foo"
	inBody, outBody := &foo{A: "a"}, `{"A":"a"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	if got, want := req.Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("NewRequest() Content-Type is %v, want %v", got, want)
	}

	if got, want := req.Header.Get("Accept"), "application/json"; got != want {
		t.Errorf("NewRequest() Accept is %v, want %v", got, want)
	}

	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}
}

func TestDo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("Request method: %v, want %v", got, want)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"A":"a"}`))
	})

	req, _ := client.NewRequest("GET", ".", nil)
	body := new(foo)
	client.Do(context.Background(), req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	resp, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}

	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}
