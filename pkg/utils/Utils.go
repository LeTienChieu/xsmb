package utils

import (
	"sort"
	"strings"
)

type ResFormatForHuman struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func CountOccurrences(list []string) map[string]int {
	countMap := make(map[string]int)
	for _, item := range list {
		countMap[item]++
	}
	return countMap
}

func SortMapByValueDesc(m map[string]int) []ResFormatForHuman {
	var sortedList []ResFormatForHuman
	for k, v := range m {
		var resultForHuman = ResFormatForHuman{Key: k, Value: v}
		sortedList = append(sortedList, resultForHuman)
	}
	sort.Slice(sortedList, func(i, j int) bool {
		return strings.Compare(sortedList[i].Key, sortedList[j].Key) > 0
	})
	return sortedList
}
