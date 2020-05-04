package main

import (
	"reflect"
	"strings"
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
				"roll": {"up": 0.5, "out": 0.5},
				"up":   {"and": 1.0},
				"and":  {"roll": 1.0},
			},
			[]string{"roll"},
		},
		{
			"keep it sweet, keep it simple, and keep your cool",
			10,
			map[string]map[string]float64{
				"keep":    {"it": 2.0 / 3, "your": 1.0 / 3},
				"it":      {"sweet,": 0.5, "simple,": 0.5},
				"sweet,":  {"keep": 1.0},
				"simple,": {"and": 1.0},
				"and":     {"keep": 1.0},
				"your":    {"cool": 1.0},
			},
			[]string{"keep"},
		},
		{
			"multiline\ncorpus",
			10,
			map[string]map[string]float64{
				"multiline": {"\n": 1.0},
				"\n":        {"corpus": 1.0},
			},
			[]string{"multiline", "corpus"},
		},
		{
			"",
			0,
			map[string]map[string]float64{},
			[]string{},
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

func TestGenerateSpeechWithNumWords(t *testing.T) {
	tests := []struct {
		corpus             string
		maxRetries         int
		numWordsToGenerate int
		numWordsWant       int
	}{
		{
			"the quick brown fox jumps over the lazy dog.\n",
			20, 42, 42,
		},
		{
			"the quick brown fox jumps over the lazy dog.\n",
			20, 0, 0,
		},
		{
			"the quick brown fox jumps over the lazy dog.\n",
			20, -1, 0,
		},
		{
			"the quick brown fox\njumps over the lazy dog.\n",
			20, 42, 42,
		},
		{
			"",
			20, 42, 0,
		},
	}
	for _, c := range tests {
		hmm := NewHMM(c.corpus, c.maxRetries)
		speech := hmm.GenerateSpeechWithNumWords(c.numWordsToGenerate)
		got := len(strings.Fields(speech))

		if got != c.numWordsWant {
			t.Errorf("GenerateSpeechWithNumWords(%d) generated speech with unexpected length: %d\n",
				c.numWordsWant, got)
		}
	}
}
