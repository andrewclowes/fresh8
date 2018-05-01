package event

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/andrewclowes/fresh8/feed"
	"github.com/andrewclowes/fresh8/importer/common"
	"github.com/andrewclowes/fresh8/store"
)

func newGetEventIdsStep(football common.FootballService) common.Step {
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

func newGetEventStep(football common.FootballService, mapEvent eventMapper) common.StepHandler {
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

func newGetMarketsStep(football common.FootballService, mapMarket marketMapper) common.StepHandler {
	return func(in interface{}) (interface{}, error) {
		event, ok := in.(*store.Event)
		if !ok {
			return nil, fmt.Errorf("invalid type for GetMarketsStep: %v", reflect.TypeOf(in))
		}

		m := make(chan *feed.Market)
		var wg sync.WaitGroup
		for _, i := range event.Markets {
			n, err := strconv.Atoi(i.ID)
			if err != nil {
				continue
			}
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				market, _, err := football.GetMarket(context.Background(), id)
				if err == nil {
					m <- market
				}
			}(n)
		}
		go func() {
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

func newStoreSendStep(eventStore common.StoreService) common.StepHandler {
	return func(in interface{}) (interface{}, error) {
		event, ok := in.(*store.Event)
		if !ok {
			return nil, fmt.Errorf("invalid type for StoreSendStep: %v", reflect.TypeOf(in))
		}

		_, err := eventStore.Create(context.Background(), event)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// NewPipeline creates a new pipeline for events
func NewPipeline(config common.ConfigProvider, football common.FootballService, store common.StoreService) (*common.Steps, error) {
	t, err := config.GetInt("jobs.event.steps.limit")
	if err != nil {
		return nil, err
	}

	eventMapper := newEventMapper(parseEventTime)
	optionMapper := newOptionMapper(parseOdds)
	marketMapper := newMarketMapper(optionMapper)

	steps := []common.StepRunner{
		newGetEventIdsStep(football),
		common.NewRateLimitedStep(newGetEventStep(football, eventMapper), t),
		common.NewRateLimitedStep(newGetMarketsStep(football, marketMapper), t),
		common.NewRateLimitedStep(newStoreSendStep(store), t),
	}
	pipeline := common.NewSteps(steps...)
	return &pipeline, nil
}
