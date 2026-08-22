// Go 1.27: ジェネリックメソッド（Generic Methods）
//
// メソッド宣言が「レシーバの型パラメータとは別に」独自の型パラメータを
// 持てるようになった。1.26 以前はこれが構文エラーで、代わりにトップレベル
// 関数（例: func MapTo[E, R any](s *Stack[E], f func(E) R) []R）として
// 書くしかなかった。1.27 からは型の名前空間の中にジェネリックな振る舞いを
// 置ける。
//
// 制約:
//   - インターフェースのメソッドはジェネリック化できない
//   - レシーバ側の型パラメータ（下の E）とメソッド側（R）は別物
//
// 実行: go run ./genericMethods
package main

import (
	"fmt"
	"strconv"
)

type Stack[E any] struct {
	items []E
}

func (s *Stack[E]) Push(e E) { s.items = append(s.items, e) }

// MapTo は「メソッド固有の型パラメータ R」を持つジェネリックメソッド。
// レシーバの E とは独立して、E から任意の R への変換を表現できる。
// これが Go 1.27 で解禁された書き方。
func (s *Stack[E]) MapTo[R any](f func(E) R) []R {
	out := make([]R, len(s.items))
	for i, e := range s.items {
		out[i] = f(e)
	}
	return out
}

func main() {
	s := &Stack[int]{}
	s.Push(1)
	s.Push(2)
	s.Push(3)

	// int -> string
	strs := s.MapTo(func(n int) string { return "#" + strconv.Itoa(n) })
	fmt.Printf("int -> string: %v\n", strs)

	// int -> bool（偶数判定）
	evens := s.MapTo(func(n int) bool { return n%2 == 0 })
	fmt.Printf("int -> bool  : %v\n", evens)
}
