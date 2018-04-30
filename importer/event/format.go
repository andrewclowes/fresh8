package event

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	timeFormats = [...]string{
		"2006-01-02:15:04:05Z",
		"2006-01-02 15:04:05Z",
	}
)

type timeParser func(raw string) (*time.Time, error)

func parseEventTime(raw string) (*time.Time, error) {
	for _, f := range timeFormats {
		if t, err := time.ParseInLocation(f, raw, time.UTC); err == nil {
			return &t, nil
		}
	}
	if i, err := strconv.ParseInt(raw, 10, 64); err == nil {
		t := time.Unix(i, 0).UTC()
		return &t, nil
	}
	return nil, fmt.Errorf("failed to parse time: %v", raw)
}

type odds struct {
	Num int
	Den int
}

type oddsParser func(raw string) (*odds, error)

func parseOdds(raw string) (*odds, error) {
	o := strings.Split(raw, "/")
	if len(o) != 2 {
		return nil, fmt.Errorf("failed to parse odds %v", raw)
	}
	nums := [2]int{}
	for i, p := range o {
		r, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("failed to parse odds %v", raw)
		}
		nums[i] = r
	}
	return &odds{Num: nums[0], Den: nums[1]}, nil
}
