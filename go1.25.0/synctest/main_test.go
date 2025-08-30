package main

import (
	"testing"
	"testing/synctest"
)

func TestRead(t *testing.T) {
	// synctestを使ってRead関数のテストを行う
	// チャネルを作成して値を送信し、Read関数を
	// 旧の場合, 最低でも10秒待機するテストがあるので10秒はかかる
	// type tests map[string]struct {
	// 	ch          chan int
	// 	expectError bool
	// }
	// testCases := tests{
	// 	"valid read": {
	// 		ch:          make(chan int, 1),
	// 		expectError: false,
	// 	},
	// 	"timeout read": {
	// 		ch:          make(chan int),
	// 		expectError: true,
	// 	},
	// }
	// for name, tc := range testCases {
	// 	t.Run(name, func(t *testing.T) {
	// 		if tc.expectError {
	// 			// チャネルに値を送信しない場合、タイムアウトが発生することを期待
	// 			_, err := Read(tc.ch)
	// 			if err == nil {
	// 				t.Fatal("expected timeout error, got nil")
	// 			}
	// 		} else {
	// 			// チャネルに値を送信して、正常に読み取れることを期待
	// 			tc.ch <- 42 // チャネルに値を送信
	// 			value, err := Read(tc.ch)
	// 			if err != nil {
	// 				t.Fatalf("unexpected error: %v", err)
	// 			}
	// 			if value != 42 {
	// 				t.Fatalf("expected value 42, got %d", value)
	// 			}
	// 		}
	// 	})
	// }

	// synctestを使うことで本来最低でも10秒かかるテストがすぐに終わる
	synctest.Test(t, func(t *testing.T) {
		ch := make(chan int)
		_, err := Read(ch)
		if err == nil {
			t.Fatal("expected timeout error, got nil")
		}
		if err.Error() == "timeout" {
			t.Log("Timeout error as expected")
		} else {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
