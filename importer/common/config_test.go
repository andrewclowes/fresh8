package common

import "testing"

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
			got, ok := p.Get(tt.Key)
			if !ok {
				t.Fatalf("get could not the key: %v", tt.Key)
			}
			if want := tt.Value; got != want {
				t.Errorf("get = %+v, want %+v", got, want)
			}
		})
	}
}
