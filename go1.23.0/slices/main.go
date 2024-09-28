package main

import (
	"fmt"
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

func main() {
	values()
}
