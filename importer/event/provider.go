package event

// import (
// 	"context"

// 	"github.com/andrewclowes/fresh8/feed"
// 	"github.com/andrewclowes/fresh8/importer/common"
// )

// // // Football defines the football service for the feed api
// // type Football interface {
// // 	ListEventIds(ctx context.Context) ([]int, *http.Response, error)
// // 	GetEvent(ctx context.Context, id int) (*feed.Event, *http.Response, error)
// // 	GetMarket(ctx context.Context, id int) (*feed.Market, *http.Response, error)
// // }

// // Provider fetches the events from the feed API
// type Provider struct {
// 	limiter  common.Limiter
// 	football Football
// }

// // Get executes the fetch of events
// func (p *Provider) Get() ([]Event, error) {
// 	ctx := context.Background()

// 	eventIds, _, err := p.football.ListEventIds(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	events, err := p.getFeedEvents(ctx, eventIds)
// 	if err != nil {
// 		return nil, err
// 	}

// 	marketIds := p.getFeedMarketIds(events)

// 	markets, err := p.getFeedMarkets(ctx, marketIds)
// 	if err != nil {
// 		return nil, err
// 	}

// }

// func (p *Provider) getFeedEvents(ctx context.Context, ids []int) ([]*feed.Event, error) {
// 	res := make(chan *feed.Event)
// 	for _, id := range ids {
// 		p.limiter.Execute(func() {
// 			e, _, err := p.football.GetEvent(ctx, id)
// 			res <- e
// 		})
// 	}
// 	p.limiter.Wait()
// 	e := make([]*feed.Event, 0)
// 	for r := range res {
// 		e = append(e, r)
// 	}
// 	return e, nil
// }

// func (p *Provider) getFeedMarkets(ctx context.Context, ids []int) ([]*feed.Market, error) {
// 	res := make(chan *feed.Market)
// 	for _, id := range ids {
// 		p.limiter.Execute(func() {
// 			m, _, err := p.football.GetMarket(ctx, id)
// 			res <- m
// 		})
// 	}
// 	p.limiter.Wait()
// 	m := make([]*feed.Market, 0)
// 	for r := range res {
// 		m = append(m, r)
// 	}
// 	return m, nil
// }

// func (p *Provider) getFeedMarketIds(events []*feed.Event) []int {
// 	marketIds := make([]int, 0)
// 	exists := map[int]bool{}
// 	for _, e := range events {
// 		if e.Markets == nil {
// 			continue
// 		}
// 		for _, m := range *e.Markets {
// 			if _, ok := exists[m]; ok {
// 				continue
// 			}
// 			exists[m] = true
// 			marketIds = append(marketIds, m)
// 		}
// 	}
// 	return marketIds
// }
