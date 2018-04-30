package event

// func TestEventMapperStep(t *testing.T) {
// 	input := &eventResources{
// 		Event: &feed.Event{
// 			ID:      apiutil.JSONNumber("1"),
// 			Name:    apiutil.String("A v B"),
// 			Time:    apiutil.String("2018-04-25T12:00:00Z"),
// 			Markets: &[]json.Number{"2"},
// 		},
// 		Markets: []*feed.Market{
// 			&feed.Market{
// 				ID:   apiutil.JSONNumber("1"),
// 				Type: apiutil.String("win"),
// 				Options: &[]feed.Option{
// 					feed.Option{
// 						ID:   apiutil.JSONNumber("1"),
// 						Name: apiutil.String("win"),
// 						Odds: apiutil.String("1/2"),
// 					},
// 				},
// 			},
// 		},
// 	}

// 	want := &Event{
// 		ID:   1,
// 		Name: "A v B",
// 		Time: time.Date(2018, 4, 25, 12, 0, 0, 0, time.UTC),
// 		Markets: []Market{
// 			Market{
// 				ID:   "1",
// 				Type: "win",
// 				Options: []Option{
// 					Option{
// 						ID:   "1",
// 						Name: "win",
// 						Num:  1,
// 						Den:  2,
// 					},
// 				},
// 			},
// 		},
// 	}

// 	s := newMapEventStep()

// 	got, err := s(input)
// 	t.Log(got)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Error("invalid results", got, want)
// 	}
// }
