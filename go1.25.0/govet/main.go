package main

import (
	"fmt"
	"sync"
)

// go vetでwait groupのエラーを検知
// go vetするとエラーになるよ
// 便利だね
func main() {
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1) //本来はこの位置にやらないとエラー
		go func(i int) {
			defer wg.Done()
			fmt.Println("i: ", i)
			// 処理
		}(i)
	}
	wg.Wait()

	// これが正解
	for i := range 5 {
		go func(i int) {
			wg.Add(1) //この位置はNG, addされる前にdoneされる,このエラーを検知できるようになった
			defer wg.Done()
			fmt.Println("i: ", i)
			// 処理
		}(i)
	}
	wg.Wait()
}
