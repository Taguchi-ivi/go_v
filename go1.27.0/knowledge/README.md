# go1.27.0 ナレッジ

実機（`go1.27.0 darwin/arm64`）で検証しながらまとめた学習ノート。
各ドキュメントは対応する `../<機能>/main.go` と対になっている。

| # | ドキュメント | 対応コード |
| --- | --- | --- |
| 00 | [概要 & 移行チェックリスト](./00-overview.md) | - |
| 01 | [ジェネリックメソッド](./01-generic-methods.md) | `../genericMethods` |
| 02 | [encoding/json/v2](./02-json-v2.md) | `../jsonV2` |
| 03 | [goroutine リークプロファイル](./03-goroutine-leak.md) | `../goroutineLeak` |
| 04 | [構造体リテラルの柔軟化](./04-struct-literal.md) | `../structLiteral` |
| 05 | [モダナイズ & 最新の書き方](./05-modernize.md) | `go fix` |
| 06 | [モダナイズ before→after カタログ](./06-modernize-examples.md) | `go fix -diff` |
| 07 | [uuid パッケージ（新規）](./07-uuid.md) | `../uuid` |
| 08 | [Minor changes（実務ピック）](./08-minor-changes.md) | `../minorChanges` |

## 実行方法

```bash
cd go1.27.0
go run ./genericMethods
go run ./jsonV2
go run ./structLiteral
go run ./goroutineLeak
go run ./modernize
go run ./uuid
go run ./minorChanges
go test ./minorChanges -v   # synctest の例
go fix -diff ./...          # モダナイズ差分をドライラン表示
```

> 前提: `GOTOOLCHAIN=auto`（このディレクトリは go.mod の `go 1.27.0` により自動で 1.27 を使う）。
