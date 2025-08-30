package main

import (
	"fmt"
	"sync"
)

func main() {
	// add, doneが不要になった

	// 旧のsync.WaitGroupはadd, doneを使ってgoroutineの完了を待機していた
	fmt.Println("old sync.WaitGroup example:")
	oldMsg := make([]string, 0, 5)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := range 5 {
		wg.Add(1) // 追加
		go func(i int) {
			defer wg.Done() // 完了
			println(i)
			mu.Lock()
			oldMsg = append(oldMsg, fmt.Sprintf("goroutine %d completed", i))
			mu.Unlock()
		}(i)
	}

	wg.Wait() // 全てのgoroutineが完了するまで待機
	fmt.Println("Old messages:", oldMsg)

	// 新しいsync.WaitGroupは新しいAPIを提供している
	// カウント漏れやdeferの忘れを防ぐ
	newMsg := make([]string, 0, 5)
	fmt.Println("new sync.WaitGroup example:")
	for i := range 5 {
		wg.Go(func() {
			println(i) // iはループの変数なので、最後の値が出力される
			mu.Lock()
			newMsg = append(newMsg, fmt.Sprintf("goroutine %d completed", i))
			mu.Unlock()
		})
	}
	wg.Wait() // 全てのgoroutineが完了するまで待機
	fmt.Println("New messages:", newMsg)
}
