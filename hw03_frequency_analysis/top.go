package hw03frequencyanalysis

import (
	"regexp"
	"sort"
)

const topCount = 10

var splitWordsRegexp *regexp.Regexp

func init() {
	splitWordsRegexp = regexp.MustCompile(`[^\s+]+`)
}

func Top10(textToAnalise string) []string {
	words := splitWordsRegexp.FindAllString(textToAnalise, -1)
	wordsCount := make(map[string]int)
	for _, v := range words {
		wordsCount[v]++
	}

	getNextKeysWithMaxValue := getNextKeysWithMaxValueClosure(wordsCount)

	top := make([]string, 0, topCount)
	for {
		keys := getNextKeysWithMaxValue()
		if keys == nil {
			break
		}

		if len(keys) > 1 {
			sort.Strings(keys)
		}

		top = append(top, keys...)
		if len(top) > topCount {
			top = top[:topCount]
			break
		}
	}

	return top
}

func getNextKeysWithMaxValueClosure(wordsMap map[string]int) func() []string {
	wordsMapCopy := make(map[string]int)
	for k, v := range wordsMap {
		wordsMapCopy[k] = v
	}

	return func() []string {
		if len(wordsMapCopy) == 0 {
			return nil
		}

		keys := getKeysWithMaxValue(wordsMapCopy)
		for _, key := range keys {
			delete(wordsMapCopy, key)
		}

		return keys
	}
}

func getKeysWithMaxValue(m map[string]int) []string {
	if len(m) == 0 {
		return nil
	}

	maxValue := getMaxValueOfMap(m)
	keysWithMaxValue := make([]string, 0, 1)
	for key, v := range m {
		if maxValue == v {
			keysWithMaxValue = append(keysWithMaxValue, key)
		}
	}

	return keysWithMaxValue
}

func getMaxValueOfMap(m map[string]int) int {
	maxValue := 0
	for _, v := range m {
		if maxValue < v {
			maxValue = v
		}
	}

	return maxValue
}
