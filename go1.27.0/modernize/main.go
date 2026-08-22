// Go 1.27 で「新しく追加 / 改名 / 削除」された go fix モダナイザだけを集めた例。
//
// 各関数は after（新しい書き方）を実際に動くコードとして書き、
// before（旧い書き方）はコメントに残している。
//
// 差分確認: go fix -diff ./modernize
// 一覧    : go tool fix help
package main

import (
	"errors"
	"fmt"
	"slices"
	"sync/atomic"
	"unsafe"
)

// 🆕 embedlit : 埋め込み(昇格)フィールドを直接初期化（1.27の構造体リテラル新機能そのもの）
type base struct{ ID int }
type wrap struct {
	base
	Tag string
}

func embedlitDemo() wrap {
	// before(〜1.26): wrap{base: base{ID: 1}, Tag: "x"}
	// after (1.27) :
	return wrap{ID: 1, Tag: "x"}
}

// 🆕 atomictypes : sync/atomic の基本型を atomic 型に置き換え
func atomictypesDemo() int64 {
	// before: var n int64; atomic.AddInt64(&n, 1); return atomic.LoadInt64(&n)
	// after :
	var n atomic.Int64
	n.Add(1)
	return n.Load()
}

// 🆕 slicesbackward : 逆順ループを slices.Backward に置き換え
func slicesbackwardDemo(s []int) []int {
	// before: for i := len(s) - 1; i >= 0; i-- { out = append(out, s[i]) }
	// after :
	out := make([]int, 0, len(s))
	for _, v := range slices.Backward(s) {
		out = append(out, v)
	}
	return out
}

// 🆕 unsafefuncs : unsafe ポインタ演算を関数呼び出し(unsafe.Add)に置き換え
func unsafefuncsDemo() int {
	arr := [3]int{10, 20, 30}
	p := unsafe.Pointer(&arr[0])
	// before: *(*int)(unsafe.Pointer(uintptr(p) + unsafe.Sizeof(arr[0])))
	// after :
	second := *(*int)(unsafe.Add(p, unsafe.Sizeof(arr[0])))
	return second
}

// 🆕 errorsastype : errors.As を errors.AsType[T] に置き換え（AsType は 1.27 の新関数）
type myErr struct{ code int }

func (e *myErr) Error() string { return fmt.Sprintf("code=%d", e.code) }

func errorsastypeDemo(err error) (int, bool) {
	// before: var t *myErr; if errors.As(err, &t) { return t.code, true }
	// after :
	if t, ok := errors.AsType[*myErr](err); ok {
		return t.code, true
	}
	return 0, false
}

// ── 改名・削除（1.27）※アナライザ名の変更なので before/after コードは無し ──
//
// 改名: waitgroup → waitgroupgo
//        （wg.Add(1)/go/wg.Done() → wg.Go(func(){...}) への変換。名前の曖昧さ回避）
// 削除: fmtappendf
//        （fmt.Sprintf(...) を []byte 化する変換。スタイル上の理由で削除）

func main() {
	fmt.Printf("embedlit       : %+v\n", embedlitDemo())
	fmt.Printf("atomictypes    : %d\n", atomictypesDemo())
	fmt.Printf("slicesbackward : %v\n", slicesbackwardDemo([]int{1, 2, 3}))
	fmt.Printf("unsafefuncs    : %d\n", unsafefuncsDemo())

	code, ok := errorsastypeDemo(fmt.Errorf("wrap: %w", &myErr{code: 42}))
	fmt.Printf("errorsastype   : code=%d ok=%v\n", code, ok)
}
