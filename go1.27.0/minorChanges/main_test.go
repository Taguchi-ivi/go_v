package main

import (
	"testing"
	"testing/synctest"
	"time"
)

// Go 1.27: testing/synctest に Sleep ヘルパーが追加。
//
// synctest のバブル内では時計が「仮想時間」で動く。synctest.Sleep(d) は
// 仮想時計を d だけ進め、かつ（自分以外の）全 goroutine が durably-blocked に
// なるまで待つ。time.Sleep + synctest.Wait を 1 行で書けるイメージ。
//
// 実時間を待たずに「2時間後」の挙動をテストできる。
//
// 実行: go test ./minorChanges -run TestSynctestSleep -v
func TestSynctestSleep(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		done := make(chan struct{})
		go func() {
			time.Sleep(2 * time.Hour) // 仮想時間で2時間ブロック
			close(done)
		}()

		// 1時間だけ進める → まだ完了していないはず
		synctest.Sleep(1 * time.Hour)
		select {
		case <-done:
			t.Fatal("2時間経過前に完了してしまった")
		default:
			t.Log("1時間経過: まだ未完了（期待通り）")
		}

		// さらに1時間 → 合計2時間 → 完了するはず
		synctest.Sleep(1 * time.Hour)
		select {
		case <-done:
			t.Log("2時間経過: 完了（期待通り）")
		default:
			t.Fatal("2時間経過したのに未完了")
		}
	})
}
