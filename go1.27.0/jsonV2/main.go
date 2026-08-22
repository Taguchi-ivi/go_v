// Go 1.27: 新しい JSON エンジン encoding/json/v2 と encoding/json/jsontext
//
// - encoding/json/v2 : 高レベル API（Marshal / Unmarshal）。デフォルトが厳格。
//     * 不正な UTF-8 を拒否
//     * オブジェクト内の重複キーを拒否（v1 は黙って後勝ち）
//     * オプションは可変長引数（...Options）で渡す
// - encoding/json/jsontext : 低レベルのトークン/整形。WithIndent などのオプションもここ。
//
// v1（encoding/json）は引き続き使える。1.27 では内部実装が v2 ベースに
// 切り替わる（GOEXPERIMENT=nojsonv2 でオプトアウト可、1.28 で削除予定）。
// v1 API はほぼ互換だが、エラーメッセージ文字列は変わりうる点に注意。
//
// ── 移行の結論（実機確認済み） ─────────────────────────────
// v2 を「自分で import しない限り」は、1.27 に上げても基本大丈夫。
//   * encoding/json（v1）のままなら意味論は不変
//       - フィールド名の大小: 区別しない（従来通り）… {"NAME"} は name に入る
//       - 重複キー: 後勝ちで受理（従来通り）
//   * 厳格化（大小区別・重複拒否・UTF-8拒否）は encoding/json/v2 を
//     import した時だけ発動する完全オプトイン機能。
//   * 1.28 で消えるのは実装フラグ GOEXPERIMENT=nojsonv2 だけ。
//     v1 の意味論は維持され、v2 への書き換えは永久に「任意」。
//   * 上げるだけで実際に変わりうるのはエラーメッセージ文字列くらい。
// 速度: v1 API のままでも v2 バックエンドの恩恵で Unmarshal は高速化する
//       （Marshal はほぼ同等）。詳細は knowledge/02-json-v2.md。
// ───────────────────────────────────────────────
//
// 実行: go run ./jsonV2
package main

import (
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"fmt"
	"os"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email,omitzero"` // ゼロ値なら出力から省く（1.24+ の omitzero タグ）
	Age   int    `json:"age"`
}

func main() {
	// 1) Marshal + 整形（jsontext.WithIndent）
	u := User{Name: "Alice", Age: 30} // Email は空 -> omitzero で消える
	b, err := json.Marshal(u, jsontext.WithIndent("  "))
	if err != nil {
		panic(err)
	}
	fmt.Println("[marshal]")
	os.Stdout.Write(b)
	fmt.Println()

	// 2) v2 の厳格さ: 重複キーはエラー
	dup := []byte(`{"name":"Bob","name":"Bob2","age":40}`)
	var u2 User
	if err := json.Unmarshal(dup, &u2); err != nil {
		fmt.Println("[dup key] rejected:", err)
	}

	// 3) 通常の Unmarshal
	var u3 User
	if err := json.Unmarshal([]byte(`{"name":"Carol","age":25}`), &u3); err != nil {
		panic(err)
	}
	fmt.Printf("[unmarshal] %+v\n", u3)
}
