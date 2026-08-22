# uuid パッケージ（新規・標準ライブラリ）

Go 1.27 で **`uuid`** がトップレベルの標準パッケージとして追加（`import "uuid"`）。
`github.com/google/uuid` などの外部依存なしで UUID 生成・解析ができる。RFC 9562 準拠。
乱数成分は暗号学的に安全な生成器を使用。

## API（`go doc uuid` 実測）

```go
type UUID [16]byte   // comparable（== が使える）

// 生成
func New() UUID       // 多くの用途に適したアルゴリズム（汎用）
func NewV4() UUID      // 完全ランダム
func NewV7() UUID      // 時刻順（先頭に Unix 時刻） ← DB主キー/ソート向き
func Nil() UUID        // 全ゼロ
func Max() UUID        // 全 FF

// 解析
func Parse(s string) (UUID, error)
func MustParse(s string) UUID

// メソッド
func (u UUID) String() string
func (u UUID) Compare(v UUID) int
func (u UUID) MarshalText() ([]byte, error)
func (u UUID) AppendText(b []byte) ([]byte, error)
func (u *UUID) UnmarshalText(b []byte) error
```

## ⚠️ `New()` は V4（ランダム）＝ V7 ではない（実機確認）

「New が一番良いやつ」と思いがちだが、**`New()` が返すのは version 4（完全ランダム）**。
version nibble を読むと確認できる:

```
New()   version = 4
NewV4() version = 4   ← New() と同じ
NewV7() version = 7
```

DB 主キーなどでソート性が欲しいなら、**明示的に `NewV7()` を呼ぶ**必要がある。

## 実務での選び方

- **DB主キー / 時系列で並べたい → `NewV7()` を明示**。時刻が先頭に来るため単調増加に近く、
  B-tree インデックスの断片化が起きにくい／`Compare` で生成順にソートされる（実機確認済み）。
- ランダム性だけで順序不要 → `NewV4()`（＝`New()`）。
- V7 のソート性: 生成順にシャッフルしても `Compare` でソートすると生成順に戻る（確認済み）。
- `UUID` は `[16]byte` かつ comparable なので、map のキーや `==` 比較がそのまま使える。
- `MarshalText`/`UnmarshalText` 実装済みなので JSON 等のテキスト表現も自然に扱える。
- 生の16バイトが欲しいときは `[16]byte(u)` に変換（`%x` はデフォルトだと文字列表現の16進になる点に注意）。

## 実行結果（`go run ./uuid`）

```
New   : 579ea0a7-bf64-40ab-9b28-10169fbf3a8b
V4    : fc23b5c4-94e7-45d3-9b7a-19b7bdaab4e9
V7(a) : 01a028ca-e190-7a4d-a59f-06f4e78b9563
V7(b) : 01a028ca-e190-7a51-8149-e0a964b34b21
V7 順序 a<=b : true
```

## 関連
- コード: `../uuid/main.go`
- [[00-overview]]
