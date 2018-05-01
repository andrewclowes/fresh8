package common

import (
	"strconv"
	"testing"
)

func TestInMemoryConfigProvider_Get(t *testing.T) {
	settingTests := []struct {
		Key   string
		Value string
	}{
		{"a", "1"},
		{"b", "2"},
		{"c", "3"},
	}
	c := make(map[string]string)
	for _, s := range settingTests {
		c[s.Key] = s.Value
	}

	p := InMemoryConfigProvider{values: c}

	for _, tt := range settingTests {
		t.Run(tt.Key, func(t *testing.T) {
			got, err := p.Get(tt.Key)
			if err != nil {
				t.Fatalf("get could not the key: %v", tt.Key)
			}
			if want := tt.Value; got != want {
				t.Errorf("get = %+v, want %+v", got, want)
			}
		})
	}
}

func TestInMemoryConfigProvider_GetInt(t *testing.T) {
	settingTests := []struct {
		Key   string
		Value int
	}{
		{"a", 1},
		{"b", 1},
		{"c", 3},
	}
	c := make(map[string]string)
	for _, s := range settingTests {
		c[s.Key] = strconv.Itoa(s.Value)
	}

	p := InMemoryConfigProvider{values: c}

	for _, tt := range settingTests {
		t.Run(tt.Key, func(t *testing.T) {
			got, err := p.GetInt(tt.Key)
			if err != nil {
				t.Fatalf("get could not the key: %v", tt.Key)
			}
			if want := tt.Value; got != want {
				t.Errorf("get = %+v, want %+v", got, want)
			}
		})
	}
}
