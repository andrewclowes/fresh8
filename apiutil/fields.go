package apiutil

import (
	"encoding/json"
	"time"
)

// Int returns a pointer to an int
func Int(x int) *int {
	return &x
}

// String returns a pointer to a string
func String(x string) *string {
	return &x
}

// Time returns a pointer to a time
func Time(x time.Time) *time.Time {
	return &x
}

// JSONNumber retuens a pointer to a time
func JSONNumber(x json.Number) *json.Number {
	return &x
}
