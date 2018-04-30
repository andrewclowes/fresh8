package common

import (
	"sync"
	"testing"
)

func TestLimiter(t *testing.T) {
	limit := 10
	n := 100

	l := NewLimiter(limit)
	m := map[int]bool{}
	lock := &sync.Mutex{}

	max := int32(0)
	for i := 0; i < n; i++ {
		x := i
		l.Execute(func() {
			lock.Lock()
			m[x] = true
			currentMax := l.GetNumInProgress()
			if currentMax >= max {
				max = currentMax
			}
			lock.Unlock()
		})
	}

	l.Wait()

	if len(m) != int(n) {
		t.Error("invalid num of results", len(m))
	}

	if max > int32(limit) {
		t.Error("invalid max", max)
	}
}
