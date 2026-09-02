package main

import (
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed data/words.txt
var rawWordList string

var wordList = loadWordList()

func loadWordList() []string {
	lines := strings.Split(rawWordList, "\n")
	var words []string
	for _, line := range lines {
		w := strings.TrimSpace(line)
		if w != "" {
			words = append(words, w)
		}
	}

	return words
}

func shuffleWords() []string {
	w := make([]string, len(wordList))
	copy(w, wordList)

	rand.Shuffle(len(w), func(i, j int) { w[i], w[j] = w[j], w[i] })
	return w
}
