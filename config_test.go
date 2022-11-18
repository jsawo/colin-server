package main

import (
	"testing"
	"time"
)

func TestGetFrequency(t *testing.T) {
	type getFrequencyCase struct {
		input    string
		expected time.Duration
	}
	cases := []getFrequencyCase{
		{"10s", time.Duration(10) * time.Second},
		{"8h", time.Duration(8) * time.Hour},
		{"3h20m", time.Duration(3)*time.Hour + time.Duration(20)*time.Minute},
	}

	for _, testCase := range cases {
		c := &Collector{Frequency: testCase.input}
		got := c.GetFrequency()
		want := testCase.expected

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	}
}
