// Go 1.27: 標準ライブラリに uuid パッケージが追加された（import "uuid"）。
//
// これまで github.com/google/uuid などの外部依存で入れていた UUID 生成が
// 標準だけで完結する。RFC 9562 準拠。乱数は暗号学的に安全な生成器を使う。
//
// ┌──────────┬──────────────────────────────┬────────────────────────────────────────┐
// │          │              V4              │                   V7                   │
// ├──────────┼──────────────────────────────┼────────────────────────────────────────┤
// │ 中身     │ 完全ランダム                 │ 先頭に Unix時刻 + ランダム             │
// ├──────────┼──────────────────────────────┼────────────────────────────────────────┤
// │ 順序     │ バラバラ                     │ 生成順にソートされる（単調増加に近い） │
// ├──────────┼──────────────────────────────┼────────────────────────────────────────┤
// │ DB主キー │ インデックスが断片化しやすい │ 断片化しにくい・範囲検索に強い         │
// ├──────────┼──────────────────────────────┼────────────────────────────────────────┤
// │ 使う場面 │ 順序不要のID                 │ DBの主キー・時系列で並べたいID         │
// └──────────┴──────────────────────────────┴────────────────────────────────────────┘
// ⚠️ 落とし穴：New() は「一番いいやつ」ではなく V4（ランダム）
// 実行: go run ./uuid
package main

import (
	"fmt"
	"sync"
	"uuid" // トップレベルの標準パッケージ
)

func main() {
	// New : 多くの用途に適したアルゴリズムで生成（汎用）
	fmt.Println("New   :", uuid.New())

	// NewV4 : 完全ランダム。従来もっとも一般的だったバージョン。
	fmt.Println("V4    :", uuid.NewV4())

	// NewV7 : 時刻順（先頭に Unix 時刻）。単調増加に近く、DBの主キーや
	//         ソートに強い。実務ではこれを選ぶ場面が多い。
	a := uuid.NewV7()
	b := uuid.NewV7()
	fmt.Println("V7(a) :", a)
	fmt.Println("V7(b) :", b)
	fmt.Println("V7 順序 a<=b :", a.Compare(b) <= 0) // 生成順にほぼ整列

	// Parse / MustParse : 文字列 -> UUID
	parsed, err := uuid.Parse(a.String())
	fmt.Printf("Parse : %v err=%v\n", parsed, err)

	// UUID は [16]byte で comparable（== が使える）
	fmt.Println("==    :", a == parsed)

	// Nil（全ゼロ）/ Max（全FF）
	fmt.Println("Nil   :", uuid.Nil())
	fmt.Println("Max   :", uuid.Max())

	// UUID は [16]byte なので、生の16バイトが必要なら配列に変換する
	raw := uuid.NewV4()
	fmt.Printf("bytes : %x (%d bytes)\n", [16]byte(raw), len(raw))

	// V7 を 10 個「並行」生成して、毎回すべて異なることを確認する。
	//
	// 【最新の書き方 / Go 1.25+】
	//   従来: wg.Add(1); go func(){ defer wg.Done(); ... }()
	//   最新: wg.Go(func(){ ... })          ← Add/Done が要らない1行
	//
	// 【ループ変数 / Go 1.22+】
	//   for i := range 10 で反復ごとに i が独立。クロージャで capture しても安全。
	//
	// 【データ競合なしの集約】
	//   各 goroutine は results の「自分の番号」だけに書く（別要素＝別メモリ）。
	//   共有スライスを resize しないので mutex 不要。順序も保てる。
	var wg sync.WaitGroup
	results := make([]uuid.UUID, 10)
	for i := range 10 {
		wg.Go(func() {
			results[i] = uuid.NewV7()
		})
	}
	wg.Wait()

	fmt.Println("--- 並行生成した V7（実行のたびに全部変わる） ---")
	for i, id := range results {
		fmt.Printf("V7[%d] : %s\n", i, id)
	}
}
