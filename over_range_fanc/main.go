package main

import "fmt"

// go run over_range_fanc/main.go

// wiki: https://tip.golang.org/wiki/RangefuncExperiment
func Backward[E any](s []E) func(func(int, E) bool) {
	return func(yield func(int, E) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(i, s[i]) {
				return
			}
		}
	}
}

// イテレータ
// Genericsもだね。
// もっとちゃんと読みたい
func main() {
	fmt.Println("Hello, World!")
	Backward([]int{1, 2, 3, 4, 5})(func(i int, v int) bool {
		fmt.Println(i, v)
		return true
	})
}
