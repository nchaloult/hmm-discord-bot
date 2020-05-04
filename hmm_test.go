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
			"keep it sweet, keep it simple, and keep your cool",
			10,
			map[string]map[string]float64{
				"keep": map[string]float64{
					"it":   2.0 / 3,
					"your": 1.0 / 3,
				},
				"it": map[string]float64{
					"sweet,":  0.5,
					"simple,": 0.5,
				},
				"sweet,":  map[string]float64{"keep": 1.0},
				"simple,": map[string]float64{"and": 1.0},
				"and":    map[string]float64{"keep": 1.0},
				"your":   map[string]float64{"cool": 1.0},
			},
			[]string{"keep"},
		},
		{
			"multiline\ncorpus",
			10,
			map[string]map[string]float64{
				"multiline": map[string]float64{ "\n": 1.0 },
				"\n": map[string]float64{ "corpus": 1.0 },
			},
			[]string{"multiline", "corpus"},
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
