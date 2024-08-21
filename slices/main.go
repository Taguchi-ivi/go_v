package main

import (
	"fmt"
	"maps"
	"slices"
)

// go run slices/main.go
// wiki: https://pkg.go.dev/slices@master

// values
func values() {

	// valueのみ
	names := []string{"Alice", "Bob", "Vera"}
	for v := range slices.Values(names) {
		fmt.Println(v)
	}
}

// sorted, keys
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

func main() {
	values()
	sortedAndKeys()
}
