package apiutil

import "time"

// Int returns a pointer to an int
func Int(x int) *int {
	return &x
}

// String returns a pointer to a string
func String(x string) *string {
	return &x
}

// Time retuens a pointer to a time
func Time(x time.Time) *time.Time {
	return &x
}
