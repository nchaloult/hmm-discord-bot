package main

import (
	"reflect"
	"testing"
)

// TestHMMCreation makes sure that NewHMM() is returning HMM structs with expected fields. In other
// words, we're effectivley testing the functionality of getWords() and buildHMMFields() in hmm.go.
func TestHMMCreation(t *testing.T) {
	tests := []struct {
		corpus         string
		maxRetries     int
		probMapWant    map[string]map[string]float64
		firstWordsWant []string
	}{
		{
			"roll up and roll out",
			20,
			map[string]map[string]float64{
				"roll": map[string]float64{
					"up":  0.5,
					"out": 0.5,
				},
				"up":  map[string]float64{"and": 1.0},
				"and": map[string]float64{"roll": 1.0},
			},
			[]string{"roll"},
		},
		{
			"",
			0,
			map[string]map[string]float64{},
			[]string{""},
		},
	}
	for _, c := range tests {
		got := NewHMM(c.corpus, c.maxRetries)
		if !reflect.DeepEqual(got.probMap, c.probMapWant) {
			t.Errorf("Unexpected probMap generated from: %q\ngot: %v\nwant: %v\n",
				c.corpus,
				got.probMap,
				c.probMapWant)
		}
		if !reflect.DeepEqual(got.firstWords, c.firstWordsWant) {
			t.Errorf("Unexpected firstWords slice generated from: %q\ngot: %v\nwant: %v\n",
				c.corpus,
				got.firstWords,
				c.firstWordsWant)
		}
		if got.maxRetries != c.maxRetries {
			t.Errorf("NewHMM() returned HMM struct with unexpected maxRetries\ngot: %d, want%d\n",
				got.maxRetries,
				c.maxRetries)
		}
	}
}
