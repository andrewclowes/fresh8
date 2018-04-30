package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// EventHandler will take the event and logs the result
func EventHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)

	log.Println(fmt.Sprintf("Request: %v", string(data)))

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
