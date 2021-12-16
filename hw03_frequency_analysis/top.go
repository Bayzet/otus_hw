package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const countTop int = 10

type Pair struct {
	World string
	Total int
}

func Top10(text string) (r []string) {
	if text == "" {
		return r
	}

	allWords := countAllWordsInText(text)
	srtWords := sortByValuesAndLexicographically(allWords)

	top := srtWords
	if len(srtWords) > countTop {
		top = srtWords[:countTop]
	}

	resultTop := make([]string, 0)

	i := 0
	for i < len(top) {
		resultTop = append(resultTop, top[i].World)
		i++
	}

	return resultTop
}

func countAllWordsInText(text string) map[string]int {
	worlds := make(map[string]int)

	wordsSlice := strings.Fields(text)
	for _, world := range wordsSlice {
		if _, ok := worlds[world]; !ok {
			worlds[world] = 0
		}

		worlds[world]++
	}

	return worlds
}

func sortByValuesAndLexicographically(words map[string]int) []Pair {
	p := make([]Pair, len(words))

	i := 0
	for k, v := range words {
		p[i] = Pair{k, v}
		i++
	}

	sort.Slice(p, func(i, j int) bool {
		switch p[i].Total == p[j].Total {
		case true:
			return p[i].World < p[j].World
		case false:
			return p[i].Total > p[j].Total
		default:
			return true
		}
	})

	return p
}
