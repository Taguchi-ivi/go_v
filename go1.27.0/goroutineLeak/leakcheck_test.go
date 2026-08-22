package main

// 実務でのリーク検査レシピ。
//
// A. CI で明示アサート          … assertNoLeaks を使う（このファイル）
// B. synctest を使うなら自動検知 … バブル終了時にブロック残があれば勝手に失敗
// C. 本番稼働サーバ             … net/http/pprof の /debug/pprof/goroutineleak
//                                （テスト不要・後述コメント参照）
//
// 実行: go test ./goroutineLeak -v

import (
	"bytes"
	"io"
	"runtime/pprof"
	"testing"
	"testing/synctest"
	"time"
)

// assertNoLeaks: goroutineleak プロファイルを走らせ、
// リークが1件でもあれば発生箇所付きで失敗させる。CIに仕込める。
func assertNoLeaks(t *testing.T) {
	t.Helper()
	p := pprof.Lookup("goroutineleak")
	p.WriteTo(io.Discard, 1) // 検出（GC）をトリガ。Count はこの後に有効
	if n := p.Count(); n > 0 {
		var buf bytes.Buffer
		p.WriteTo(&buf, 1)
		t.Fatalf("goroutine leak detected: %d\n%s", n, buf.String())
	}
}

// A. 正常系: 後始末する処理はリーク0でパスする。
func TestNoLeak(t *testing.T) {
	done := make(chan int)
	go func() { done <- 1 }() // 受信されるので終了する
	<-done
	time.Sleep(5 * time.Millisecond)
	assertNoLeaks(t)
}

// B. synctest を使えばリーク検知は「自動」。
// バブル終了時にブロック中の goroutine が残ると synctest.Test が
// deadlock で失敗する（= 専用アサートを書かなくても落ちる）。
func TestSynctestAutoLeak(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch := make(chan int, 1)
		go func() { ch <- 1 }()
		<-ch // 全goroutineが終了 -> リークなし -> パス
	})
	// ↑ ここで go func(){ <-make(chan int) }() のように放置すると、
	//   synctest.Test が "deadlock: ... blocked goroutines remain" で失敗する。
}

// C. 本番サーバでは（テスト不要）:
//     import _ "net/http/pprof"
//     go http.ListenAndServe("localhost:6060", nil)
//   稼働中に:
//     curl http://localhost:6060/debug/pprof/goroutineleak?debug=1
//   でリーク中の goroutine とスタックがそのまま得られる。
