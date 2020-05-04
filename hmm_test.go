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
		expectedErr    error
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
			nil,
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
			nil,
		},
		{
			"multiline\ncorpus",
			10,
			map[string]map[string]float64{
				"multiline": {"\n": 1.0},
				"\n":        {"corpus": 1.0},
			},
			[]string{"multiline", "corpus"},
			nil,
		},
		{
			"",
			0,
			map[string]map[string]float64{},
			[]string{},
			ErrEmtpyCorpus,
		},
		{
			"foo",
			-1,
			map[string]map[string]float64{},
			[]string{},
			ErrNegMaxRetries,
		},
	}
	for _, c := range tests {
		got, err := NewHMM(c.corpus, c.maxRetries)
		if (err == nil && c.expectedErr != nil) || (err != nil && c.expectedErr == nil) ||
			(c.expectedErr != nil && err != c.expectedErr) {
			t.Fatalf("Unexpected error. got: %v\nwant: %v\n", err, c.expectedErr)
		}
		if c.expectedErr != nil && err == c.expectedErr {
			// We expected a specific error, and we got that error. We're done with this test case.
			continue
		}

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
			t.Errorf("Unexpected maxRetries. got: %d, want%d\n",
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
	}
	for _, c := range tests {
		hmm, _ := NewHMM(c.corpus, c.maxRetries)
		speech := hmm.GenerateSpeechWithNumWords(c.numWordsToGenerate)
		got := len(strings.Fields(speech))

		if got != c.numWordsWant {
			t.Errorf("Unexpected speech length. got: %d, want: %d\n",
				got, c.numWordsWant)
		}
	}
}

func TestGenerateSpeechBeginningWithWord(t *testing.T) {
	tests := []struct {
		corpus        string
		maxRetries    int
		firstWordWant string
	}{
		{
			"the quick brown fox jumps over the lazy dog.\n",
			42, "lazy",
		},
		{
			"the quick brown fox\njumps over the lazy dog.\n",
			42, "lazy",
		},
		{
			"the quick brown fox jumps over the lazy dog.\n",
			42, "foo",
		},
		{
			"the quick brown fox\njumps over the lazy dog.\n",
			42, "foo",
		},
	}
	for _, c := range tests {
		hmm, _ := NewHMM(c.corpus, c.maxRetries)
		speech := hmm.GenerateSpeechBeginningWithWord(c.firstWordWant)

		// Get first word from speech
		firstSpaceIndex := strings.Index(speech, " ")
		if firstSpaceIndex == -1 {
			firstSpaceIndex = len(speech)
		}
		firstWord := speech[:firstSpaceIndex]

		if firstWord != c.firstWordWant {
			t.Errorf("Unexpected first word in generated speech. got: %q, want: %q\n",
				firstWord, c.firstWordWant)
		}
	}
}

func TestGenerateSpeechBeginningWithWordAndWithNumWords(t *testing.T) {
	tests := []struct {
		corpus             string
		maxRetries         int
		firstWordWant      string
		numWordsToGenerate int
		numWordsWant       int
	}{
		{
			"the quick brown fox jumps over the lazy dog.\n",
			42, "lazy", 42, 42,
		},
		{
			"the quick brown fox\njumps over the lazy dog.\n",
			42, "lazy", 42, 42,
		},
		{
			"the quick brown fox jumps over the lazy dog.\n",
			42, "foo", 42, 42,
		},
		{
			"the quick brown fox\njumps over the lazy dog.\n",
			42, "foo", 42, 42,
		},
		{
			"the quick brown fox jumps over the lazy dog.\n",
			42, "lazy", 0, 0,
		},
		{
			"the quick brown fox\njumps over the lazy dog.\n",
			42, "lazy", 0, 0,
		},
		{
			"the quick brown fox jumps over the lazy dog.\n",
			42, "foo", 0, 0,
		},
		{
			"the quick brown fox\njumps over the lazy dog.\n",
			42, "foo", 0, 0,
		},
		{
			"the quick brown fox jumps over the lazy dog.\n",
			42, "lazy", -1, 0,
		},
		{
			"the quick brown fox\njumps over the lazy dog.\n",
			42, "lazy", -1, 0,
		},
		{
			"the quick brown fox jumps over the lazy dog.\n",
			42, "foo", -1, 0,
		},
		{
			"the quick brown fox\njumps over the lazy dog.\n",
			42, "foo", -1, 0,
		},
	}
	for _, c := range tests {
		hmm, _ := NewHMM(c.corpus, c.maxRetries)
		speech := hmm.GenerateSpeechBeginningWithWordAndWithNumWords(c.firstWordWant, c.numWordsToGenerate)

		got := len(strings.Fields(speech))
		if got != c.numWordsWant {
			t.Errorf("Unexpected speech length. got %d, want: %d\n",
				got, c.numWordsWant)
		}
		if c.numWordsWant == 0 {
			// We're expecting an empty string, and we got it. We're done with this test case.
			continue
		}

		// Get first word from speech
		firstSpaceIndex := strings.Index(speech, " ")
		if firstSpaceIndex == -1 {
			firstSpaceIndex = len(speech)
		}
		firstWord := speech[:firstSpaceIndex]

		if firstWord != c.firstWordWant {
			t.Errorf("Unexpected first word in generated speech. got: %q, want: %q\n",
				firstWord, c.firstWordWant)
		}
	}
}
