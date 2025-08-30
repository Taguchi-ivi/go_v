package main

import (
	"fmt"
	"reflect"
)

func main() {
	// anyを型にキャストできるかを短い文章で対応できる(変換後の値, 変換できたかを返す)
	// とても便利だね。
	v := reflect.ValueOf(1)
	i, ok := reflect.TypeAssert[int](v) // 1 true
	fmt.Printf("i: %v, ok: %v\n", i, ok)
	s, ok := reflect.TypeAssert[string](v) //  false
	fmt.Printf("s: %v, ok: %v\n", s, ok)

	// any
	var a any
	a = "こんにちは"
	vv := reflect.ValueOf(a)
	ai, ok := reflect.TypeAssert[int](vv) // 1 true
	fmt.Printf("ai: %v, ok: %v\n", ai, ok)
	as, ok := reflect.TypeAssert[string](vv) //  false
	fmt.Printf("s: %v, ok: %v\n", as, ok)
}
