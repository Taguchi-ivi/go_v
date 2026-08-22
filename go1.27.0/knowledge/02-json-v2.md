# encoding/json/v2 & encoding/json/jsontext

## パッケージ構成

| import | 役割 |
| --- | --- |
| `encoding/json/v2` | 高レベル API。`Marshal` / `Unmarshal`。オプションは可変長 `...Options` |
| `encoding/json/jsontext` | 低レベルのトークン/整形。`WithIndent` などの整形オプションもここ |
| `encoding/json`（v1） | 従来 API。1.27 では内部が v2 実装になる（`GOEXPERIMENT=nojsonv2` で戻せる） |

慣例として `import json "encoding/json/v2"` と別名を付けると v1 と混同しにくい。

## 入出力 API（4系統）

| 関数 | 対象 | 用途 |
| --- | --- | --- |
| `Marshal` / `Unmarshal` | `[]byte` | 従来同様 |
| `MarshalWrite` / `UnmarshalRead` | `io.Writer` / `io.Reader` | **ストリーミング**（大きなデータ・HTTP body） |
| `MarshalEncode` / `UnmarshalDecode` | `jsontext.Encoder` / `Decoder` | 低レベル制御と組み合わせ |

すべて末尾に `opts ...Options` を取る。整形など構文系は `jsontext` 側のオプション。

## v2 の「厳格デフォルト」（v1 からの挙動変化）

| 挙動 | v1 | **v2 デフォルト** | 戻すオプション |
| --- | --- | --- | --- |
| フィールド名の大小 | case-**insensitive** | **case-sensitive**（`"NAME"` は `name` に入らない） | `MatchCaseInsensitiveNames(true)` |
| 重複キー | 後勝ちで受理 | **拒否**（エラー） | － |
| 不正 UTF-8 | 置換 | **拒否** | － |
| 未知フィールド | 無視 | 無視（同じ） | `RejectUnknownMembers(true)` で拒否に |
| ゼロ値省略 | `omitempty` | `json:",omitzero"` 推奨 | － |

> **これらの厳格デフォルトが効くのは `encoding/json/v2` を明示 import した時だけ。**
> `encoding/json`（v1）を使い続ける限り、1.27でも **v1 の意味論（大小区別しない・重複は後勝ち）が保たれる**（実機確認済み）。

### ⚠️ よくある誤解（重要）

| 誤解 | 実際 |
| --- | --- |
| 「1.27に上げると大小区別で弾かれる」 | ❌ v1 API のままなら従来通り。厳格化は **v2 を import した時だけ** |
| 「1.28で v2 に移行が強制される」 | ❌ 1.28 で消えるのは `GOEXPERIMENT=nojsonv2`（**実装**スイッチ）だけ。v1 の**意味論は維持**。v2 移行は永久に任意 |
| 「大小不一致はエラーになる」 | ❌ v2 でも**エラーにならず黙ってフィールドが空**になる（err=nil）。静かに危ない |

実機確認（1.27・v1 API）:
```
{"NAME":"Alice","AGE":30} -> {Name:Alice Age:30}  // v1 は大小区別せず拾う
{"name":"a","name":"b"}   -> {Name:b}             // 重複は後勝ちで受理
```

1.27 に上げるだけで実際に変わりうるのは **エラーメッセージの文字列**程度。

## 速度：v2 に移行しなくても速くなる（実測）

1.27 では v1 API `encoding/json` も v2 バックエンドで動くため、**コード無改変・v1 のまま**で
Unmarshal が速くなる。**同一の v1 コード**を 1.26 / 1.27 で回した比較（同一マシン, `-count=3`）:

| ベンチ | Go 1.26 | Go 1.27 | 差 |
| --- | --- | --- | --- |
| Unmarshal 時間 | ~6030 ns/op | **~3290 ns/op** | **約1.8倍速（-45%）** |
| Unmarshal メモリ | 1032 B/op | **688 B/op** | -33% |
| Unmarshal alloc | 28 | **15** | -46% |
| Marshal 時間 | ~1249 ns/op | ~1419 ns/op | やや遅い（+14%, 誤差〜微減） |
| Marshal メモリ | 320 B / 1 alloc | 320 B / 1 alloc | 同じ |

- **結論: Unmarshal は「上げるだけ」で大きく速くなる。v2 API への移行は速度目的では不要。**
- Marshal はほぼ横ばい（マイクロベンチでは僅かに遅く出ることも）。
- gRPC は protobuf（`google.golang.org/protobuf`）を使い `encoding/json` を通らないので、**この高速化・挙動変化はいずれも無関係**。影響するのは REST/JSON ハンドラや設定ファイルのパース等のみ。

## 主なオプション（`func(bool) Options` 形式で Marshal/Unmarshal に渡す）

`Deterministic` / `MatchCaseInsensitiveNames` / `OmitZeroStructFields` /
`RejectUnknownMembers` / `StringifyNumbers` / `FormatNilSliceAsNull` / `FormatNilMapAsNull` /
`WithMarshalers` / `WithUnmarshalers`。複数は `JoinOptions(...)` でまとめられる。

## 呼び出し側カスタム（型を触らずに変換を差し替え）

v1 では型に `MarshalJSON` を実装するしかなかったが、v2 は**呼び出し側**で型ごとの変換を注入できる。

```go
b, _ := json.Marshal(ev, json.WithMarshalers(
    json.MarshalFunc(func(t time.Time) ([]byte, error) {
        return []byte(`"` + t.Format("2006/01/02") + `"`), nil
    }),
))
// => {"at":"2026/08/22"}   ... Event 型は一切変更していない
```

外部パッケージの型や、文脈ごとに書式を変えたい時に効く（`UnmarshalFunc` も同様）。

## 使い方（`go run ./jsonV2` の実測）

```go
import (
    "encoding/json/jsontext"
    json "encoding/json/v2"
)

b, _ := json.Marshal(u, jsontext.WithIndent("  ")) // 整形はオプションで渡す
err := json.Unmarshal(data, &u, /* opts... */)
```

出力:

```
[marshal]
{
  "name": "Alice",
  "age": 30
}
[dup key] rejected: jsontext: duplicate object member name "name"
[unmarshal] {Name:Carol Email: Age:25}
```

## 移行時の注意（→ [[00-overview]] チェックリスト #2）

- **1.27に上げるだけ（v1 API のまま）なら意味論は不変**。実際に変わりうるのは
  **エラーメッセージ文字列**程度なので、エラー文字列に依存したテストだけ確認すればよい。
- 厳格挙動（大小区別・重複拒否・UTF-8拒否）は **`encoding/json/v2` を自分で import した時のみ**。
  よって「重複キー後勝ち」等に依存したコードが壊れるのは **v2 に移行した場合だけ**。
- `GOEXPERIMENT=nojsonv2` は v1 API の**内部実装**を旧版に固定するフラグ（意味論スイッチではない）。
  1.28 で削除予定＝以後 v1 API は常に v2 バックエンドで動くが、**v1 の意味論は維持**され、コード書き換えは不要。

## 関連

- コード: `../jsonV2/main.go`
- go1.25.0 の `encoding/` の続きに相当
