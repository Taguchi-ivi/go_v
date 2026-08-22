# モダナイズ（modernize）と最新の書き方

## 「modernize パッケージ」とは

もともと gopls（Go 言語サーバ）内の `modernize` アナライザ群として提供され、
エディタの Quick Fix で「古い書き方 → 新しい書き方」に変換してくれる機能。
**Go 1.27 でこれが `go fix` に統合され、コマンドラインから使えるようになった。**

- `go vet` … 「疑わしい構造」を報告（バグ検出寄り）
- `go fix` … 「改善の機会」を報告し**修正を適用**（モダナイズ寄り）

```bash
go tool fix help          # 登録アナライザ一覧
go fix -diff ./...        # 変換内容を diff で確認（適用しない・ドライラン）
go fix ./...              # 実際に適用（-fix 相当。デフォルトで全アナライザ実行）
go fix -rangeint ./...    # 特定のアナライザだけ
go fix -any=false ./...   # 特定のアナライザだけ無効化
```

## 1.27 で使える主なモダナイザ（`go tool fix help` の実データ）

| アナライザ | 変換内容 | 新しい書き方 |
| --- | --- | --- |
| `any` | `interface{}` → `any` | `var x any` |
| `rangeint` | `for i:=0;i<n;i++` → range over int | `for i := range n` |
| `minmax` | if/else → `min` / `max` 組み込み関数 | `x = max(a, b)` |
| `slicescontains` | ループ → `slices.Contains` / `ContainsFunc` | `slices.Contains(s, v)` |
| `slicessort` | `sort.Slice` → `slices.Sort` | `slices.Sort(s)` |
| `slicesbackward` | 逆順ループ → `slices.Backward` | `for i, v := range slices.Backward(s)` |
| `mapsloop` | 明示ループ → `maps` パッケージ | `maps.Copy(dst, src)` |
| `stringsseq` | `range strings.Split/Fields` → `SplitSeq`/`FieldsSeq` | 配列確保を回避 |
| `stringscut` | `strings.Index` 系 → `strings.Cut` | `strings.Cut(s, sep)` |
| `stringscutprefix` | `HasPrefix`+`TrimPrefix` → `CutPrefix` | `strings.CutPrefix(s, p)` |
| `stringsbuilder` | `+=` 連結 → `strings.Builder` | |
| `waitgroupgo` | `wg.Add(1)`/`go`/`wg.Done()` → `wg.Go` | `wg.Go(func(){ ... })`（1.25+） |
| `omitzero` | `omitempty` → `omitzero`（struct フィールド） | `json:",omitzero"` |
| `embedlit` | 埋め込みフィールド参照を簡素化 | `User{ID: 1}`（→ [[04-struct-literal]]） |
| `atomictypes` | `sync/atomic` の基本型 → atomic 型 | `atomic.Int64` |
| `newexpr` | 1.26 の `new(expr)` を活用 | `new(x)` |
| `errorsastype` | `errors.As` → `errors.AsType[T]` | `errors.AsType[*MyErr](err)` |
| `reflecttypefor` | `reflect.TypeOf(x)` → `TypeFor[T]()` | `reflect.TypeFor[T]()` |
| `testingcontext` | `context.WithCancel` → `t.Context`（テスト内） | `t.Context()` |
| `forvar` | 冗長なループ変数再宣言を削除 | （1.22 のループ変数セマンティクスに追随） |
| `rangeint`/`stditerators` 等 | イテレータAPI（`Len`/`At` → range）へ | |

> 注: `fmtappendf` はスタイル上の理由で削除、`waitgroup` は曖昧回避のため `waitgroupgo` に改名された。

## 実測: `go fix -diff` の出力例

古いスタイルのコードに対して:

```diff
-	for i := 0; i < 5; i++ {
+	for i := range 5 {
 ...
-	var x interface{} = 1
+	var x any = 1
```

- `rangeint` と `any` は自動変換された。
- `wg.Add(1)`/`go`/`wg.Done()` のループ内パターンは**保守的に見送られる**ことがある
  （変換が安全と判断できない形は触らない）。モダナイザは「安全に適用できる修正」を原則とする。

## 「最新の書き方」チェックリスト（1.25〜1.27 の潮流）

- `interface{}` は使わず **`any`**。
- カウンタループは **`for i := range n`**（1.22+）。
- 最大/最小は if/else でなく組み込み **`min` / `max`**（1.21+）。
- スライス/マップ操作は自作ループより **`slices` / `maps` パッケージ**。
- 文字列分割の反復は **`strings.SplitSeq` / `FieldsSeq`**（イテレータ、確保を回避）。
- 並行の待ち合わせは **`wg.Go(func(){...})`**（1.25+）。
- JSON のゼロ値省略は **`json:",omitzero"`**（`omitempty` より意図が明確）。
- エラー型取り出しは **`errors.AsType[T](err)`**。
- reflect の型取得は **`reflect.TypeFor[T]()`**。
- 迷ったら `go fix -diff ./...` を回して、標準の推奨形を機械的に確認する。

## 運用のコツ

- **`go fix -diff ./...` をまず回して差分を眺める** → 納得したものだけ `go fix -<name>` で個別適用、が安全。
- CI に `go fix -diff` を仕込んで「モダナイズ差分ゼロ」を保つ運用も可能（ただしチームの合意が要る）。

## 関連

- **[[06-modernize-examples]]** … 実際に `go fix` が変換した before → after カタログ
- [[04-struct-literal]]（embedlit） / [[00-overview]]
