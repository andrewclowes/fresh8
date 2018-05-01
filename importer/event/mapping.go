package event

import (
	"errors"
	"fmt"
	"time"

	"github.com/andrewclowes/fresh8/feed"
	"github.com/andrewclowes/fresh8/store"
)

type errorCollector []error

func (m *errorCollector) Collect(err error) {
	*m = append(*m, err)
}

func (m *errorCollector) Error() string {
	err := "mapping errors:\n"
	for i, e := range *m {
		err += fmt.Sprintf("\terror %d: %s\n", i, e.Error())
	}
	return err
}

type eventMapper func(*feed.Event) (*store.Event, error)

func newEventMapper(parseTime timeParser) eventMapper {
	return func(feedEvent *feed.Event) (*store.Event, error) {
		errs := new(errorCollector)
		if feedEvent == nil {
			errs.Collect(errors.New("event is nil"))
		}
		if feedEvent.ID == nil {
			errs.Collect(errors.New("event.ID is nil"))
		}
		if feedEvent.Name == nil {
			errs.Collect(errors.New("event.Name is nil"))
		}
		var t *time.Time
		var err error
		if feedEvent.Time == nil {
			errs.Collect(errors.New("event.Time is nil"))
		} else {
			t, err = parseTime(*feedEvent.Time)
			if err != nil {
				errs.Collect(err)
			}
		}
		markets := []store.Market{}
		if feedEvent.Markets == nil {
			errs.Collect(errors.New("event.Markets is nil"))
		} else {
			for _, m := range *feedEvent.Markets {
				markets = append(markets, store.Market{ID: m.String()})
			}
		}
		if len(*errs) > 0 {
			return nil, errs
		}
		event := store.Event{
			ID:      feedEvent.ID.String(),
			Name:    *feedEvent.Name,
			Time:    *t,
			Markets: markets,
		}
		return &event, nil
	}
}

type marketMapper func(*feed.Market) (*store.Market, error)

func newMarketMapper(mapOption optionMapper) marketMapper {
	return func(feedMarket *feed.Market) (*store.Market, error) {
		errs := new(errorCollector)
		if feedMarket == nil {
			errs.Collect(errors.New("market is nil"))
		}
		if feedMarket.ID == nil {
			errs.Collect(errors.New("market.ID is nil"))
		}
		if feedMarket.Type == nil {
			errs.Collect(errors.New("market.Type is nil"))
		}
		opts := []store.Option{}
		if feedMarket.Options == nil {
			errs.Collect(errors.New("market.Options is nil"))
		} else {
			for _, o := range *feedMarket.Options {
				opt, err := mapOption(&o)
				if err != nil {
					errs.Collect(err)
				} else {
					opts = append(opts, *opt)
				}
			}
		}
		if len(*errs) > 0 {
			return nil, errs
		}
		market := store.Market{
			ID:      feedMarket.ID.String(),
			Type:    *feedMarket.Type,
			Options: opts,
		}
		return &market, nil
	}
}

type optionMapper func(*feed.Option) (*store.Option, error)

func newOptionMapper(parseOdds oddsParser) optionMapper {
	return func(feedOption *feed.Option) (*store.Option, error) {
		errs := new(errorCollector)
		if feedOption == nil {
			errs.Collect(errors.New("option is nil"))
		}
		if feedOption.ID == nil {
			errs.Collect(errors.New("option.ID is nil"))
		}
		if feedOption.Name == nil {
			errs.Collect(errors.New("option.Name is nil"))
		}
		var o *odds
		var err error
		if feedOption.Odds == nil {
			errs.Collect(errors.New("option.Odds is nil"))
		} else {
			o, err = parseOdds(*feedOption.Odds)
			if err != nil {
				errs.Collect(err)
			}
		}
		if len(*errs) > 0 {
			return nil, errs
		}
		option := store.Option{
			ID:   feedOption.ID.String(),
			Name: *feedOption.Name,
			Num:  o.Num,
			Den:  o.Den,
		}
		return &option, nil
	}
}
