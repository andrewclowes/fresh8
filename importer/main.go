package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/andrewclowes/fresh8/feed"
	"github.com/andrewclowes/fresh8/importer/common"
	"github.com/andrewclowes/fresh8/importer/event"
	"github.com/andrewclowes/fresh8/store"
)

const (
	clientTimeout = 10
)

func main() {
	logger := common.NewLogger()
	config, err := common.NewConfigProvider()
	if err != nil {
		log.Fatalln(err)
	}

	timeout, err := config.GetInt("services.client.timeout")
	if err != nil {
		timeout = clientTimeout
	}
	netClient := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	feed, err := createFeedClient(config, netClient)
	if err != nil {
		log.Fatalln(err)
	}
	store, err := createStoreClient(netClient)
	if err != nil {
		log.Fatalln(err)
	}

	event, err := event.NewPipeline(config, feed.Football, store.Event)
	if err != nil {
		log.Fatalln(err)
	}
	eventJob := common.NewPipelineJob(event, logger)

	runner := common.NewJobRunner()
	runner.Run(eventJob)
}

func createFeedClient(config common.ConfigProvider, client *http.Client) (*feed.Client, error) {
	h, err := config.Get("services.client.football.host")
	if err != nil {
		return nil, err
	}
	p, err := config.Get("services.client.football.port")
	if err != nil {
		return nil, err
	}
	f, err := feed.NewClient(fmt.Sprintf("http://%v:%v", h, p), client)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func createStoreClient(client *http.Client) (*store.Client, error) {
	h := os.Getenv("STORE_ADDR")
	if h == "" {
		return nil, fmt.Errorf("environment variable STORE_ADDR not present")
	}
	f, err := store.NewClient(fmt.Sprintf("http://%v", h), client)
	if err != nil {
		return nil, err
	}
	return f, nil
}
