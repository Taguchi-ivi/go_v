package main

import (
	"fmt"
)

// go run slices/main.go
// wiki: https://pkg.go.dev/slices@master

func main() {
	names := []string{"Alice", "Bob", "Vera"}
	// for v := range slices.Values(names) {
	// 	fmt.Println(v)
	// }
	for i, v := range names {
		fmt.Println(i, v)
	}
}
