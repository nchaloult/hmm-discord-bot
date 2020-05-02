package main

import (
	"math/rand"
	"strings"
	"time"
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

	// The max number of times that speech generation is allowed to restart.See GenerateSpeech() for
	// more details.
	maxRetries int
}

// NewHMM returns a new HMM with fields populated based on the provided corpus file.
func NewHMM(corpus string, maxRetries int) *HMM {
	words := getWords(corpus)
	probMap, firstWords := buildHMMFields(words)

	return &HMM{
		probMap:    probMap,
		firstWords: firstWords,
		maxRetries: maxRetries,
	}
}

// GenerateSpeech returns a piece of generated text. After it finishes generating a sentence, a
// counter called: retries is incremented by a random number between 1 and 3. Once retries is
// greater than or equal to maxRetries, all of the sentences that were generated are returned.
func (h *HMM) GenerateSpeech() string {
	var speech []string
	retries := 0

	// Seeding the pseudo-random number generator is necessary:
	// https://stackoverflow.com/questions/33994677/pick-a-random-value-from-a-go-slice
	n := len(h.firstWords)
	rand.Seed(time.Now().Unix())
	curWord := h.firstWords[rand.Intn(n)]

	for retries < h.maxRetries {
		speech = append(speech, curWord)
		curWord = getNextWord(curWord, h.probMap)

		if curWord == "\n" {
			retries += rand.Intn(2) + 1 // Generate int in range: [1, 3]
		}
	}

	output := strings.Join(speech, " ")
	return strings.TrimSpace(output)
}

// GenerateSpeechWithNumWords returns a piece of generated text with the provided number of words.
func (h *HMM) GenerateSpeechWithNumWords(numWords int) string {
	var speech []string

	n := len(h.firstWords)
	rand.Seed(time.Now().Unix())
	curWord := h.firstWords[rand.Intn(n)]

	for i := 0; i < numWords; i++ {
		speech = append(speech, curWord)
		curWord = getNextWord(curWord, h.probMap)
	}

	output := strings.Join(speech, " ")
	return strings.TrimSpace(output)
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

// getNextWord consults the provided probMap to pick the next word that should follow the provided
// curWord like a hidden Markov model would.
func getNextWord(curWord string, probMap map[string]map[string]float64) string {
	if successorProbs, ok := probMap[curWord]; ok {
		cur := 0.0
		for successor, prob := range successorProbs {
			cur += prob
			if rand.Float64() <= cur {
				return successor
			}
		}
	}

	// If the provided curWord isn't in probMap, then pick a starting word at random from the list
	// of probMap's keys.
	n := len(probMap)
	probMapKeps := make([]string, n)
	i := 0
	for key := range probMap {
		probMapKeps[i] = key
		i++
	}
	rand.Seed(time.Now().Unix())
	return probMapKeps[rand.Intn(n)]
}
