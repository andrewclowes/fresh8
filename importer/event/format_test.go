package event

import (
	"reflect"
	"testing"
	"time"
)

func TestEventTimeParser(t *testing.T) {
	want := time.Date(2018, 4, 25, 12, 0, 0, 0, time.UTC)
	timeTests := []struct {
		input string
		want  time.Time
	}{
		{"2018-04-25:12:00:00Z", want},
		{"2018-04-25 12:00:00Z", want},
		{"1524657600", want},
	}

	for _, tt := range timeTests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseEventTime(tt.input)
			if err != nil {
				t.Fatalf("parseEventTime returned error: %v", err)
			}
			if *got != tt.want {
				t.Errorf("parseEventTime = %+v, want %+v", *got, &tt.want)
			}
		})
	}
}

func TestEventOddsParser(t *testing.T) {
	oddsTests := []struct {
		input string
		want  odds
	}{
		{"1/2", odds{Num: 1, Den: 2}},
		{"2/5", odds{Num: 2, Den: 5}},
		{"5/1", odds{Num: 5, Den: 1}},
	}

	for _, tt := range oddsTests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseOdds(tt.input)
			if err != nil {
				t.Fatalf("parseOdds returned error: %v", err)
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("parseOdds = %+v, want %+v", *got, &tt.want)
			}
		})
	}
}
