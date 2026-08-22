// Go 1.27 の「Minor changes to the library」から、実務で使いそうなものだけを抜粋。
//
// 実行: go run ./minorChanges
// synctest の例は main_test.go 側（go test ./minorChanges）。
package main

import (
	"fmt"
	rand "math/rand/v2"
	"net/url"
	"strings"
)

func main() {
	// ── strings.CutLast / bytes.CutLast ─────────────────────────────
	// 既存の Cut は「最初の」区切りで分割。CutLast は「最後の」区切りで分割。
	// パスと拡張子、ファイル名など「末尾側で切りたい」場面に便利。
	path := "src/pkg/util/helper.go"
	dir, file, ok := strings.CutLast(path, "/")
	fmt.Printf("[CutLast] dir=%q file=%q ok=%v\n", dir, file, ok)

	name, ext, _ := strings.CutLast("archive.tar.gz", ".")
	fmt.Printf("[CutLast] name=%q ext=%q\n", name, ext) // name="archive.tar" ext="gz"

	// ── net/url URL.Clone / Values.Clone ────────────────────────────
	// ディープコピー。ミドルウェアで「元を壊さずにクエリを足す」等で重宝する。
	base, _ := url.Parse("https://api.example.com/v1/items?page=1")
	req := base.Clone()
	q := req.Query()
	q.Set("page", "2")
	req.RawQuery = q.Encode()
	fmt.Printf("[Clone] base=%s\n", base) // 元は無傷
	fmt.Printf("[Clone] req =%s\n", req)

	// ── math/rand/v2 Rand.N（ジェネリックメソッド：1.27新機能の標準ライブラリ実例）──
	// 任意の整数型で [0,n) の乱数。型変換なしで書ける。
	r := rand.New(rand.NewPCG(42, 1024))
	fmt.Printf("[Rand.N] int32=%d uint=%d\n", r.N(int32(100)), r.N(uint(1_000_000)))
}
