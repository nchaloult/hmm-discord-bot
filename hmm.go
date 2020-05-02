package main

import (
	"strings"
)

// HMM generates pieces of text with the same vocabulary and sentence structure as a corpus file. A
// hidden Markov model is used to generate content.
type HMM struct {
	// Collection of words in the corpus and the ratios of appearance for all words that follow.
	//
	// Ex: if the corpus looks like: "roll up and roll out", then prob looks like:
	// { "roll": { "up": 0.5,
	//             "out": 0.5 },
	//   "up": { "and": 1.0 },
	//   "and": { "roll": 1.0 },
	//   "out": { "\n": 1.0 } }
	probMap map[string]map[string]float64

	// List of words that appear at the beginning of new lines in the corpus.
	firstWords []string
}

// NewHMM returns a new HMM with fields populated based on the provided corpus file.
func NewHMM(corpus string) *HMM {
	words := getWords(corpus)
	probMap, firstWords := buildHMMFields(words)

	return &HMM{
		probMap:    probMap,
		firstWords: firstWords,
	}
}

// getWords performs input sanitization on the provided string, and splits it up into a slice of
// words.
//
// Can't use strings.Fields() instead of strings.Split() because that func nukes newline chars.
func getWords(corpus string) []string {
	paddedNewlines := strings.ReplaceAll(corpus, "\n", " \n ")
	lowercase := strings.ToLower(paddedNewlines)
	return strings.Split(lowercase, " ")
}

// buildHMMFields builds an HMM's prob field and firstWords field from a provided slice of words.
func buildHMMFields(words []string) (map[string]map[string]float64, []string) {
	freqMap := make(map[string]map[string]int)
	firstWords := []string{words[0]}
	// Populate freqMap and firstWords
	for i := 0; i < len(words)-1; i++ {
		cur := words[i]
		successor := words[i+1]

		if _, ok := freqMap[cur]; !ok {
			freqMap[cur] = make(map[string]int)
			freqMap[cur][successor] = 1
		} else if _, ok := freqMap[cur][successor]; !ok {
			freqMap[cur][successor] = 1
		} else {
			freqMap[cur][successor]++
		}

		if cur == "\n" {
			firstWords = append(firstWords, successor)
		}
	}

	probMap := make(map[string]map[string]float64)
	// Populate prob
	for cur, successors := range freqMap {
		probMap[cur] = make(map[string]float64)
		numCurOccurrences := 0
		for _, freq := range successors {
			numCurOccurrences += freq
		}
		for successor := range successors {
			probMap[cur][successor] = float64(successors[successor]) / float64(numCurOccurrences)
		}
	}

	return probMap, firstWords
}
