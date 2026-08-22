# 構造体リテラルの柔軟化（埋め込みフィールドの直接初期化）

## 何が変わったか

埋め込み（embedded）型の**昇格フィールドを、ネストせずフィールド名で直接初期化**できるようになった。

```go
type Base struct { ID int; Name string }
type Meta struct { Tag string }
type User struct {
    Base
    Meta
    Age int
}
```

| バージョン | 書き方 |
| --- | --- |
| 1.26 以前 | `User{Base: Base{ID: 1, Name: "Alice"}, Age: 30}`（`ID:` を直接書くと `unknown field ID`） |
| 1.27 以降 | `User{ID: 1, Name: "Alice", Tag: "vip", Age: 30}` ← `Base.` / `Meta.` を省略できる |

従来のネスト記法（`Base: Base{...}`）も引き続き有効。混在も可。

## 実測（`go run ./structLiteral`）

```
{Base:{ID:1 Name:Alice} Meta:{Tag:vip} Age:30}
{Base:{ID:2 Name:Bob} Meta:{Tag:} Age:40}
```

## モダナイザ連携

`go fix` の **`embedlit`**（"simplify references to embedded fields in composite literals"）が
旧ネスト記法 → 直接記法へ自動変換してくれる。→ [[05-modernize]]

## 境界（実機で確認）

| ケース | 可否 | 備考 |
| --- | --- | --- |
| 2段以上のネスト埋め込みの昇格フィールド | ✅ OK | `C{X:1}`（C←B←A の X）も直接書ける |
| 従来のネスト記法との混在 | ✅ OK | `User{Base: Base{...}, Age: 30}` も併用可 |
| 名前が衝突する昇格フィールド | ❌ error | `unknown field ID in struct literal of type M`（曖昧なので従来通り修飾が必要） |
| ポインタ埋め込み(`*A`)越しの昇格フィールド | ❌ error | `invalid implicit pointer indirection to reach X`（暗黙のアロケーションが必要になるため） |

## ハマりどころ

- IDE(gopls) が 1.27 未対応だと `unknown field ID in struct literal [MissingLitField]` を赤線表示する。
  → gopls を 1.27 でビルドし直せば解消（本リポジトリでは対応済み）。
- 名前衝突・ポインタ埋め込みは上表の通り**直接記法が使えない**ので、その場合は従来のネスト記法で書く。

## 関連

- コード: `../structLiteral/main.go`
