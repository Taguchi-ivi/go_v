// Go 1.27: 構造体リテラルの柔軟化（埋め込みフィールドの直接初期化）
//
// 埋め込み（embedded）型の昇格フィールドを、ネストせずに
// フィールド名で直接初期化できるようになった。
//
//	1.26 以前:  User{Base: Base{ID: 1, Name: "Alice"}, Age: 30}
//	1.27 以降:  User{ID: 1, Name: "Alice", Age: 30}   // Base. を書かなくてよい
//
// これは go fix の `embedlit` モダナイザ
// （"simplify references to embedded fields in composite literals"）で
// 旧 -> 新へ自動変換できる。
//
// 実行: go run ./structLiteral
package main

import "fmt"

type Base struct {
	ID   int
	Name string
}

type Meta struct {
	Tag string
}

type User struct {
	Base // 埋め込み
	Meta // 埋め込み
	Age  int
}

func main() {
	// 昇格フィールド ID / Name / Tag を直接指定できる
	u := User{
		ID:   1,
		Name: "Alice",
		Tag:  "vip",
		Age:  30,
	}
	fmt.Printf("%+v\n", u)

	// 従来のネスト記法ももちろん有効
	u2 := User{
		Base: Base{ID: 2, Name: "Bob"},
		Age:  40,
	}
	fmt.Printf("%+v\n", u2)
}
