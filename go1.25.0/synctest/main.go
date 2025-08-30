package main

import (
	"errors"
	"fmt"
	"time"
)

func Read(ch <-chan int) (int, error) {
	select {
	case v := <-ch:
		return v, nil // チャネルから値を読み取る
	case <-time.After(10 * time.Second):
		return 0, errors.New("timeout") // タイムアウト
	}
}

func main() {
	// ここではRead関数をテストするためのコードを書くことができます
	// 例えば、チャネルを作成して値を送信し、Read関数を呼び出すなど
	ch := make(chan int, 1)
	ch <- 42 // チャネルに値を送信, 即座に送るとtimeoutにはならない. ここをコメントアウトするとtimeoutになる

	value, err := Read(ch)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Read value:", value) // 読み取った値を表示
}
