package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	sortedAndKeys()
}

func sortedAndKeys() {
	m := map[string]int{
		"a": 1,
		"d": 2,
		"c": 3,
		"b": 4,
	}
	sorted := slices.Sorted(maps.Keys(m))
	for k, v := range sorted {
		fmt.Println(k, v)
	}
}
