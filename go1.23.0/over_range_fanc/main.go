package main

import (
	"fmt"
	"iter"
)

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

// iter.Seqで内容を表示する,今までとあまり変わらないやり方
func PrintAll[V any](seq iter.Seq[V]) {
	for v := range seq {
		fmt.Println(v)
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

	// https://zenn.dev/mattn/articles/641f1d86fffdc9
	// iterには二つの型が定義されている
	// type (
	// 	Seq[V any]     func(yield func(V) bool)
	// 	Seq2[K, V any] func(yield func(K, V) bool)
	// )

	// iter.Seqとiter.Seq2の違いはforの使い方。 sliceをforで回すときはSeqを使う、mapをforで回すときはSeq2を使う

}
