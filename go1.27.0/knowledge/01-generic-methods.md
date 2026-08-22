# ジェネリックメソッド（Generic Methods）

## 何が変わったか

メソッドが**レシーバの型パラメータとは独立した、自前の型パラメータ**を宣言できるようになった。

```go
func (s *Stack[E]) MapTo[R any](f func(E) R) []R { ... }
//              ^^^ レシーバ側    ^^^ メソッド側（1.27 で解禁）
```

- 1.26 以前: `method must have no type parameters` の構文エラー。
  代替としてトップレベル関数 `func MapTo[E, R any](s *Stack[E], ...)` にするしかなかった。
- 1.27: 型の名前空間の中にジェネリックな振る舞いを置ける。

## 制約・境界（1.27 実機で検証）

| 項目 | 挙動 | 確認したエラー/結果 |
| --- | --- | --- |
| インターフェースに generic method | **不可** | `interface method must have no type parameters` |
| リフレクションから見えるか | **見えない** | `reflect.Type.NumMethod()` に含まれない（通常メソッドのみ） |
| 直接呼び出し（型推論） | OK | `b.MapTo(func(int)string{...})` は引数から R を推論 |
| method value（instantiation無し） | **不可** | `cannot use generic function b.MapTo without instantiation` |
| method value（instantiation明示） | OK | `mv := b.MapTo[string]` |
| method expression | OK | `me := Box[int].MapTo[string]` |

- レシーバの型パラメータ（`E`）とメソッドの型パラメータ（`R`）は別物。
- **なぜインターフェース不可か**: ジェネリックメソッドは「型引数を与えて初めて具体的な関数」になるため、
  メソッドセット（＝動的ディスパッチ表）に載せられない。よって interface に入れられず、reflect にも現れない。
  → 「**静的にしか呼べない**（コンパイル時に型が決まる）」機能だと理解するとよい。
- 標準ライブラリの実例: `math/rand/v2` の `Rand.N[Int intType](n Int) Int`（→ [[08-minor-changes]]）。

## 実行結果（`go run ./genericMethods`）

```
int -> string: [#1 #2 #3]
int -> bool  : [false true false]
```

## ハマりどころ

- **IDE(gopls)が 1.27 未対応だと `method must have no type parameters` を赤線表示する。**
  これはエディタ側の解析が古いだけで、`go run` / `go build`（1.27）は通る。
  gopls が 1.27 対応版に上がれば解消する。

## 関連

- [[00-overview]] の言語機能一覧
- コード: `../genericMethods/main.go`
