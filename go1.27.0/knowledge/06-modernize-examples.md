# モダナイズ before → after カタログ（新しく書けるようになった形）

`go fix -diff ./...` に実際のコードを食わせて採取した**本物の変換結果**。
（Go 1.27 / `go1.27.0 darwin/arm64` で検証。各見出しの `name` は `go tool fix help` のアナライザ名）

> 使い方: `go fix -diff ./...`（差分プレビュー）→ `go fix ./...`（適用）／個別は `go fix -embedlit ./...`

---

## ✅ 検証済み（実際に `go fix` が変換したもの）

### `any` — `interface{}` → `any`
```go
// before
func f() interface{} { return nil }
// after
func f() any { return nil }
```

### `rangeint` — 3節for → for range int
```go
// before
for i := 0; i < 10; i++ { sum += i }
// after
for i := range 10 { sum += i }
```

### `minmax` — if/else → min / max
```go
// before
var m int
if a > b { m = a } else { m = b }
// after
m = max(a, b)
```

### `slicescontains` — ループ → slices.Contains
```go
// before
for _, x := range s { if x == v { return true } }
return false
// after
return slices.Contains(s, v)
```

### `slicessort` — sort.Slice → slices.Sort
```go
// before
sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
// after
slices.Sort(s)
```

### `slicesbackward` — 逆順ループ → slices.Backward
```go
// before
for i := len(s) - 1; i >= 0; i-- { sum += s[i] }
// after
for _, v := range slices.Backward(s) { sum += v }
```

### `mapsloop` — コピー用ループ → maps.Copy
```go
// before
for k, v := range src { dst[k] = v }
// after
maps.Copy(dst, src)
```

### `stringsbuilder` — += 連結 → strings.Builder
```go
// before
s := ""
for _, x := range xs { s += x }
return s
// after
var s strings.Builder
for _, x := range xs { s.WriteString(x) }
return s.String()
```

### `stringscut` — strings.Index → strings.Cut
```go
// before
i := strings.Index(s, "=")
if i < 0 { return s, "" }
return s[:i], s[i+1:]
// after
before, after, ok := strings.Cut(s, "=")
if !ok { return s, "" }
return before, after
```

### `stringscutprefix` — HasPrefix + TrimPrefix → CutPrefix
```go
// before
if strings.HasPrefix(s, "go") { return strings.TrimPrefix(s, "go"), true }
// after
if after, ok := strings.CutPrefix(s, "go"); ok { return after, true }
```

### `stringsseq` — range Split → SplitSeq（配列確保を回避）
```go
// before
for range strings.Split(doc, "\n") { n++ }
// after
for range strings.SplitSeq(doc, "\n") { n++ }
```

### `waitgroupgo` — Add/go/Done → wg.Go（Go 1.25+）
```go
// before
wg.Add(1)
go func() { defer wg.Done() /* ... */ }()
// after
wg.Go(func() { /* ... */ })
```

### `atomictypes` — sync/atomic の基本型 → atomic 型
```go
// before
var n int64
atomic.AddInt64(&n, 1)
return atomic.LoadInt64(&n)
// after
var n atomic.Int64
n.Add(1)
return n.Load()
```

### `errorsastype` — errors.As → errors.AsType[T]（Go 1.27）
```go
// before
var t *netErr
if errors.As(err, &t) { return t.msg }
// after
if t, ok := errors.AsType[*netErr](err); ok { return t.msg }
```

### `reflecttypefor` — reflect.TypeOf(値) → reflect.TypeFor[T]()
```go
// before
return reflect.TypeOf("")
// after
return reflect.TypeFor[string]()
```
> 注: 引数が「型が静的に分かる」場合のみ発火（変数だと見送られることがある）。

### `embedlit` — 埋め込みフィールドのネスト literal → 直接（Go 1.27）
```go
// before
return Wrap{Base: Base{ID: 1}, Tag: "x"}
// after
return Wrap{ID: 1, Tag: "x"}
```
> → 構造体リテラルの新機能そのもの（[[04-struct-literal]]）

### `hostport` — fmt.Sprintf("%s:%d") → net.JoinHostPort
```go
// before
net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
// after
net.Dial("tcp", net.JoinHostPort(host, fmt.Sprintf("%d", port)))
```
> IPv6 でも壊れないアドレス組み立てになる（`[::1]:80` を正しく生成）。

### `forvar` — 冗長なループ変数の再宣言を削除（Go 1.22+ のループ変数仕様）
```go
// before
for _, v := range xs { v := v; use(v) }
// after
for _, v := range xs { use(v) }
```

### `omitzero` — struct フィールドの無意味な omitempty を整理
```go
// before
Nested Inner `json:"nested,omitempty"` // struct 型に omitempty は効かない
// after
Nested Inner `json:"nested"`
```
> struct 型フィールドで「ゼロ値なら省く」をしたいなら `json:",omitzero"` を使う（v1/v2 共通、Go 1.24+）。

---

## 📋 その他の登録アナライザ（文脈依存で今回は未発火 / 参考）

| name | 変換内容 |
| --- | --- |
| `newexpr` | Go 1.26 の `new(expr)` を活用（`p := new(int); *p = x` 系） |
| `unsafefuncs` | unsafe ポインタ演算 → 関数呼び出しへ |
| `stditerators` | `Len`/`At` スタイルAPI → イテレータ（range）へ |
| `testingcontext` | テスト内 `context.WithCancel` → `t.Context()` |
| `inline` | `//go:fix inline` ディレクティブに基づくインライン展開 |
| `plusbuild` / `buildtag` | 旧 `// +build` の除去・`//go:build` 整合チェック |

> これらは「安全に適用できる」と判断された箇所でのみ変換される。発火しなくてもコードが悪いわけではない。

## 運用のおすすめ

1. `go fix -diff ./...` で差分を眺める（適用しないドライラン）。
2. 納得したものだけ個別適用: `go fix -embedlit -rangeint ./...`。
3. まとめて適用: `go fix ./...`。
4. CI に `go fix -diff` を仕込んで「モダナイズ差分ゼロ」を保つ運用も可（チーム合意前提）。

## 関連
- [[05-modernize]]（アナライザ一覧と最新の書き方チェックリスト）
- [[04-struct-literal]] / [[00-overview]]
