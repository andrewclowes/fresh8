package common

import (
	"reflect"
	"testing"
)

func TestPipeline(t *testing.T) {
	input := make(chan interface{})

	add3 := adder(3)
	mult2 := multiplier(2)

	pipeline := NewSteps(add3, mult2)
	out, _ := pipeline.Run(input)

	go func() {
		for _, i := range []int{1, 2, 3} {
			input <- i
		}
		close(input)
	}()

	got := []int{}
	for o := range out {
		r := o.(int)
		got = append(got, r)
	}

	if len(got) != 3 {
		t.Error("invalid num of results", len(got))
	}

	want := []int{8, 10, 12}
	if !reflect.DeepEqual(got, want) {
		t.Error("invalid results", got, want)
	}
}

func adder(x int) Step {
	return Step(func(in <-chan interface{}, out chan interface{}, errc <-chan error) {
		for m := range in {
			n := m.(int)
			out <- (int(n) + x)
		}
	})
}

func multiplier(x int) Step {
	return Step(func(in <-chan interface{}, out chan interface{}, errc <-chan error) {
		for m := range in {
			n := m.(int)
			out <- (int(n) * x)
		}
	})
}
