package event

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/andrewclowes/fresh8/feed"
	"github.com/andrewclowes/fresh8/importer/common"
	"github.com/andrewclowes/fresh8/store"
)

// FootballService defines the football service for the feed api
type footballService interface {
	ListEventIds(ctx context.Context) ([]json.Number, *http.Response, error)
	GetEvent(ctx context.Context, id int) (*feed.Event, *http.Response, error)
	GetMarket(ctx context.Context, id int) (*feed.Market, *http.Response, error)
}

// storeService defines the event service for the store api
type storeService interface {
	Create(ctx context.Context, event *store.Event) (*http.Response, error)
}

type eventResources struct {
	Event   *feed.Event
	Markets []*feed.Market
}

func newGetEventIdsStep(football footballService) common.Step {
	return common.Step(func(in <-chan interface{}, out chan interface{}) {
		ids, _, err := football.ListEventIds(context.Background())
		if err != nil {
			return
		}
		for _, e := range ids {
			out <- e
		}
	})
}

func newGetEventStep(football footballService, mapEvent eventMapper) common.StepHandler {
	return func(in interface{}) (interface{}, error) {
		n, ok := in.(json.Number)
		if !ok {
			return nil, fmt.Errorf("invalid type for GetEventStep: %v", reflect.TypeOf(in))
		}
		id, err := n.Int64()
		if err != nil {
			return nil, fmt.Errorf("failed conversion to int: %v", n)
		}
		event, _, err := football.GetEvent(context.Background(), int(id))
		if err != nil {
			return nil, err
		}
		e, _ := mapEvent(event)
		return e, nil
	}
}

func newGetMarketsStep(football footballService, mapMarket marketMapper) common.StepHandler {
	return func(in interface{}) (interface{}, error) {
		event, ok := in.(*store.Event)
		if !ok {
			return nil, fmt.Errorf("invalid type for GetMarketsStep: %v", reflect.TypeOf(in))
		}
		m := make(chan *feed.Market)
		go func() {
			var wg sync.WaitGroup
			for _, i := range event.Markets {
				l, err := strconv.Atoi(i.ID)
				if err != nil {
					continue
				}
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					market, _, err := football.GetMarket(context.Background(), id)
					fmt.Println(err)
					if err == nil {
						m <- market
						fmt.Println(market.ID)
					}
				}(l)
			}
			wg.Wait()
			close(m)
		}()
		markets := []store.Market{}
		for market := range m {
			n, _ := mapMarket(market)
			markets = append(markets, *n)
		}
		event.Markets = markets
		return event, nil
	}
}

func newStoreSendStep(eventStore storeService) common.StepHandler {
	return func(in interface{}) (interface{}, error) {
		event, ok := in.(*store.Event)
		if !ok {
			return nil, fmt.Errorf("invalid type for StoreSendStep: %v", reflect.TypeOf(in))
		}
		_, err := eventStore.Create(context.Background(), event)
		if err != nil {
			return nil, err
		}
		return true, nil
	}
}

// NewPipeline creates a new pipeline for events
func NewPipeline() *common.Steps {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}
	football, _ := feed.NewClient("http://localhost:8000", netClient)
	eventStore, _ := store.NewClient("http://localhost:8001", netClient)

	eventMapper := newEventMapper(parseEventTime)
	optionMapper := newOptionMapper(parseOdds)
	marketMapper := newMarketMapper(optionMapper)

	steps := []common.StepRunner{
		newGetEventIdsStep(football.Football),
		common.NewRateLimitedStep(newGetEventStep(football.Football, eventMapper), 10),
		common.NewRateLimitedStep(newGetMarketsStep(football.Football, marketMapper), 10),
		common.NewRateLimitedStep(newStoreSendStep(eventStore.Event), 10),
	}
	pipeline := common.NewSteps(steps...)
	return &pipeline
}
