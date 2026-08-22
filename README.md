# Go バージョン別学習リポジトリ

Goの各バージョンの新機能を学習・テストするためのリポジトリです。

## 📚 概要

このリポジトリは、Goの各バージョンで追加された新機能や改善点を体系的に学習するために作成されました。各バージョンごとにディレクトリを分けて、実際にコードを動かしながら理解を深めることを目的としています。

## 🗂️ プロジェクト構成

```
/
├── go1.23.0/          # Go 1.23.0 の学習・実験
│   ├── maps/          # maps パッケージの新機能
│   ├── slices/        # slices パッケージの新機能
│   ├── unique/        # unique パッケージ
│   ├── struct/        # struct の改善
│   └── over_range_fanc/  # range over func の機能
├── go1.25.0/          # Go 1.25.0 の学習・実験
│   ├── net_http/      # net/http の改善
│   ├── refrect/       # reflect パッケージの変更
│   ├── syncWaitGroup/ # sync.WaitGroup の改善
│   └── synctest/      # sync テストユーティリティ
└── go1.27.0/          # Go 1.27.0 の学習・実験（実機検証済み）
    ├── knowledge/         # 検証ノート（00〜08 + 索引）
    ├── genericMethods/    # ジェネリックメソッド
    ├── jsonV2/            # encoding/json/v2
    ├── goroutineLeak/     # goroutine リークプロファイル
    ├── structLiteral/     # 構造体リテラル（埋め込み直接init）
    ├── modernize/         # go fix モダナイズ（1.27新規分）
    ├── uuid/              # uuid パッケージ（新規）
    └── minorChanges/      # strings.CutLast / url.Clone / rand.N / synctest.Sleep
```

> ⚠️ **go1.24 / go1.26 は未学習**（ディレクトリ無し）。特に **1.26** は Green Tea GC の
> デフォルト化など重要な変更があるので、次の学習候補（下記「学習内容」参照）。

## 🚀 使い方

### 例の実行
```bash
# 特定の例を実行
cd go1.23.0/maps
go run main.go

# ルートから実行
go run ./go1.23.0/maps/main.go
```

### モジュール管理
```bash
# バージョンディレクトリ内で依存関係を整理
cd go1.23.0/
go mod tidy

# 依存関係の検証
go mod verify
```

## 📖 学習内容

### Go 1.23.0
- **maps パッケージ**: マップ操作の新機能
- **slices パッケージ**: スライス操作の新機能
- **unique パッケージ**: 一意性の保証
- **struct 改善**: 構造体関連の機能向上
- **range over func**: 関数をrangeで回せる新機能
- **go command**: `go mod tidy diff` の追加
- **Timer/Ticker**: GC の改善
- **go linkname**: 廃止方向の機能

### Go 1.25.0
- **net/http**: HTTP関連の改善
- **reflect**: リフレクション機能の変更
- **sync.WaitGroup**: 同期処理の改善
- **synctest**: テスト用同期ユーティリティ
- **encoding/json v2**: JSON処理の新バージョン

### Go 1.24.0（未学習）
- ジェネリック型エイリアス、`os.Root`、`weak` パッケージ、Swiss Table マップ 等（TODO）

### Go 1.26.0（未学習・次の学習候補）
- **Green Tea GC がデフォルト化**（GCオーバーヘッド 10〜40%削減）
- `new(expr)` 構文 ほか（TODO）
- ※ 1.27 の学習中に関連トピックとして触れているため `go1.27.0/knowledge/00-overview.md` も参照

### Go 1.27.0（学習済み・実機検証済み）
- **ジェネリックメソッド**: メソッドが独自の型パラメータを持てる（言語の目玉）
- **encoding/json/v2**: 新JSONエンジン（厳格デフォルト・v1のまま高速化・完全オプトイン）
- **goroutine リークプロファイル**: `pprof` で永久ブロック goroutine を検出（GC到達可能性で判定）
- **構造体リテラル**: 埋め込み昇格フィールドを直接初期化
- **go fix モダナイズ**: `embedlit` / `atomictypes` / `slicesbackward` / `unsafefuncs` / `errorsastype` 追加
- **uuid パッケージ**（新規, `import "uuid"`）
- **Minor changes**: `strings.CutLast` / `url.Clone` / `rand.N` / `synctest.Sleep`
- **メモリ/GC**: サイズ特化アロケータ（小割り当て最大30%高速）/ Green Tea GC 恒久化
- **移行**: 最小 macOS 13 / bzr廃止。詳細は `go1.27.0/README.md` と `go1.27.0/knowledge/`

## 📝 参考資料

### Go 1.23.0
- [Release Notes](https://tip.golang.org/doc/go1.23)
- [Go Blog](https://go.dev/blog/go1.23)
- [Future Tech Blog - slices/maps](https://future-architect.github.io/articles/20240719a/)
- [Future Tech Blog - range over func](https://future-architect.github.io/articles/20240718a/)
- [網羅的解説記事](https://zenn.dev/koya_iwamura/articles/7e7482c7222e37)
- [リリースパーティ動画](https://www.youtube.com/watch?v=SNdToIxRn20)

### Go 1.25.0
- [Release Notes](https://tip.golang.org/doc/go1.25)
- [リリース解説動画](https://www.youtube.com/watch?v=_RX5tiq0Mu8)
- [Zenn解説記事](https://zenn.dev/koya_iwamura/articles/ea2cf191cdcb2a)
- [Future Tech Blog - Go 1.25特集](https://future-architect.github.io/tags/Go1-25/)

### Go 1.27.0
- [Release Notes](https://go.dev/doc/go1.27)
- [Go Blog](https://go.dev/blog/go1.27)
- [Ready for Go 1.27 on Day One (JetBrains)](https://blog.jetbrains.com/go/2026/08/20/ready-for-go-1-27-on-day-one/)
- [Linuxiac まとめ](https://linuxiac.com/go-1-27-released-with-generic-methods-json-v2-and-faster-memory-allocation/)
- [VictoriaMetrics interactive tour](https://victoriametrics.com/blog/go-1-27/)

## 🛠️ 開発環境

- **Go Playground**: https://go.dev/play/
- **各バージョンのGo**: 対応するGoバージョンが必要

## 💡 学習のポイント

- 各バージョンは独立したgo.modを持つため、バージョン固有の機能を安全に試せます
- 実際にコードを動かして動作を確認することを重視しています
- 各ディレクトリのREADME.mdに詳細な学習メモがあります

## 🤝 貢献

このリポジトリは個人の学習目的で作成されていますが、改善提案やバグ報告は歓迎します。

## このREADMEについて

このREADMEはCLAUDEによって記載されています。変更しても変えない可能性があるので注意が必要です by taguchi-ivi