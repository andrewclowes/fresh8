package main

import (
	"net/http"
	"time"

	"github.com/andrewclowes/fresh8/feed"
	"github.com/andrewclowes/fresh8/importer/common"
	"github.com/andrewclowes/fresh8/importer/event"
	"github.com/andrewclowes/fresh8/store"
)

func main() {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}
	feed, _ := feed.NewClient("http://localhost:8000", netClient)
	store, _ := store.NewClient("http://localhost:8001", netClient)

	s := event.NewPipeline(feed.Football, store.Event)
	p := common.NewPipelineJob(s)
	p.Run()
}
