package main

// unique package
// https://pkg.go.dev/unique@master
// func main() {
// ものすごいでかい構造体を比較する時に使うかな？？

// メモリを節約できるっぽい(https://future-architect.github.io/articles/20240719a/)
// 特定のシナリオにおいては、パフォーマンスの最適化ができそう
// one := unique.Make[int](1)
// }

// NOTE: リリースパーティより
// interningというかなり攻めたメモリ使用量の削減手法をGoに実装したもの
// 使い道
// 同じ値をいろんな箇所で何度も使う場合
// 特に値がstringや容量が大きい構造体の場合。

// mattnさん
// https://github.com/mattn/go-unique-example/blob/main/main.go
// メモリ使用量の比較をしている

import (
	"log"
	"runtime"
	"unique"
)

var ss = []string{
	"fooooooooooooooooooooo",
	"baaaaaaaaaaaaaaaaaaaar",
	"baaaaaaaaaaaaaaaaaaaaz",
	"zooooooooooooooooooooo",
	"ohoooooooooooooooooooo",
}

const N = 10000000

// 現在のメモリ使用量を取得
func getAlloc() uint64 {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// uniqueパッケージを使わない場合
func test_without_unique_package() {
	before := getAlloc()
	a := make([]string, N)
	for i := 0; i < len(a); i++ {
		a[i] = ss[i%len(ss)]
	}
	log.Printf("test_without_unique_package: %v allocated", getAlloc()-before)
}

// uniqueパッケージを使う場合
func test_with_unique_package() {
	before := getAlloc()
	a := make([]unique.Handle[string], N)
	for i := 0; i < len(a); i++ {
		a[i] = unique.Make(ss[i%len(ss)])
	}
	log.Printf("test_with_unique_package: %v allocated", getAlloc()-before)
}

// 処理の比較
func main() {
	test_without_unique_package()
	runtime.GC()
	test_with_unique_package()
}

// 結果, 断然メモリ使用量が少ない
// test_without_unique_package: 21088 allocated
// test_with_unique_package: 1440 allocated
