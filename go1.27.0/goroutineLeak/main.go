// Go 1.27: goroutine リークプロファイル（goroutineleak）
//
// runtime/pprof に "goroutineleak" プロファイルが正式追加。
// チャネル / sync.Mutex / sync.Cond などの同期プリミティブで
// 「永久にブロックして解放されない」goroutine を検出する。
//
// ★ 判定の核心 = GC の「到達可能性」
//
//	単にブロックしているだけではリークにならない。
//	ブロック対象（チャネル等）が「もうどこからも到達できない」＝
//	誰も二度と signal できない状態のときだけリークと判定する。
//	（メモリの GC と同じ発想を goroutine に適用したもの）
//
// ポイント:
//   - pprof.Lookup("goroutineleak") で取得
//   - WriteTo が検出のために GC サイクルを走らせる（コスト高＝常用でなく診断用）
//     そのため Count() は WriteTo 前だと 0 のことがある
//   - net/http/pprof 経由なら /debug/pprof/goroutineleak
//   - 完全検出は原理的に不可能だが、大きなクラスのリークを拾える
//
// 実行: go run ./goroutineLeak
package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

// reachable は main から到達可能なグローバルチャネル。
// これを待つ goroutine は「まだ signal されうる」のでリークではない。
var reachable = make(chan int)

// realLeak: このgoroutine以外から到達できないチャネルを待つ = 本当のリーク。
func realLeak() {
	ch := make(chan int)
	<-ch
}

// blockedButReachable: ブロックしているが、待ち対象(reachable)は到達可能。
// → リークとは判定されない（誤検知しない）。
func blockedButReachable() {
	<-reachable
}

func main() {
	go realLeak()            // ← これだけがリーク
	go realLeak()            // ← これもリーク
	go blockedButReachable() // ← ブロック中だがリークではない
	time.Sleep(50 * time.Millisecond)

	p := pprof.Lookup("goroutineleak")
	if p == nil {
		fmt.Println("goroutineleak プロファイルが見つかりません（Go 1.27+ が必要）")
		return
	}

	// 3つブロック中だが、リークとして出るのは realLeak の2件だけ。
	fmt.Println("--- goroutineleak profile (debug=1) ---")
	p.WriteTo(os.Stdout, 1)
}

// ── 実務で使うなら（テストは必須ではない・二段構えが実用的） ──────────────
//
// C. 本番稼働サーバで観測【テスト不要・一番手軽】
//      import _ "net/http/pprof"
//      go http.ListenAndServe("localhost:6060", nil)
//    稼働中に叩くだけでリーク箇所が出る:
//      curl http://localhost:6060/debug/pprof/goroutineleak?debug=1
//    監視系にスクレイプさせてもよい。
//
// A. CI で明示アサート【予防の本命・ヘルパー1個】
//      p := pprof.Lookup("goroutineleak")
//      p.WriteTo(io.Discard, 1)                        // 検出(GC)をトリガ
//      if n := p.Count(); n > 0 { t.Fatalf("leaked %d", n) }  // 実測: Count=2
//    → goroutineLeak/leakcheck_test.go に動くテンプレあり。
//
// B. synctest を使うなら自動検知【専用アサート不要】
//    synctest.Test のバブル終了時にブロック中 goroutine が残ると
//      panic: deadlock: main bubble goroutine has exited but blocked goroutines remain
//    で勝手に失敗する。
//
// まとめ: 実地調査は C（本番でエンドポイント）／作り込み防止は A・B を CI。
// ───────────────────────────────────────────────
