package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/event", EventHandler)

	log.Println("server running on :8001")
	if err := http.ListenAndServe(":8001", r); err != nil {
		panic(err)
	}
}
