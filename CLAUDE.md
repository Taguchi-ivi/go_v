# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 最重要ルール・新しいルールの追加プロセス

ユーザーから今回限りではなく、常に対応が必要だと思われる指示を受けた場合

1. 「これを標準のルールにしますか？」と質問する。
2. YESの回答を得た場合、CLAUDE.mdに追加ルールとして記載する
3. 以降は標準ルールとして常に適用する

このプロセスによって、プロジェクトのルールを継続的に改善していきます。

追加ルール一覧
- 必ず日本語で回答してください
- このレポジトリは基本的にGoの最新versionの勉強、テストを目的としています。何か依頼をするときは最新の情報を取得して回答するようにしてください

## Project Structure

This is a Go learning repository organized by Go versions, containing study materials and example implementations for different Go language versions.

```
/
├── go1.23.0/          # Go 1.23.0 studies and examples
│   ├── maps/          # maps package examples  
│   ├── slices/        # slices package examples
│   ├── unique/        # unique package examples
│   ├── struct/        # struct improvements
│   └── over_range_fanc/  # range over func examples
└── go1.25.0/          # Go 1.25.0 studies and examples (in development)
    ├── net_http/      # net/http improvements
    ├── refrect/       # reflect package changes
    ├── syncWaitGroup/ # sync.WaitGroup improvements
    └── synctest/      # sync testing utilities
```

Each version directory contains:
- `README.md` - Study notes and reference links in Japanese
- `go.mod` - Go module definition for that version
- Individual example directories with `main.go` files

## Common Development Tasks

### Running Examples
```bash
# Navigate to specific example directory
cd go1.23.0/maps
go run main.go

# Or run from root
go run ./go1.23.0/maps/main.go
```

### Module Management
```bash
# From version directory (e.g., go1.23.0/)
go mod tidy

# Check dependencies
go mod verify
```

### Building
```bash
# Build specific example
cd go1.23.0/maps
go build .

# Build all examples in a version
find go1.23.0 -name "*.go" -execdir go build . \;
```

## Architecture Notes

- Each Go version is isolated in its own directory with separate `go.mod` files
- Examples demonstrate specific features introduced in each Go version
- go1.23.0 focuses on: maps, slices, unique, struct improvements, range over functions
- go1.25.0 focuses on: net/http improvements, reflect changes, sync.WaitGroup, sync testing
- All examples use simple `main.go` files for easy execution and learning
- Documentation is primarily in Japanese reflecting the learning context