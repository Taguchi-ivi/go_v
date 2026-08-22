
## v.1.27.0 の勉強

Go 1.27（2026年8月リリース）の学習メモとサンプル。
詳しい検証ノートは [`knowledge/`](./knowledge/README.md) にまとめている。

> **既存プロジェクトの移行メモ（実機検証済み）**
> - **JSON v2 は完全オプトイン**。`encoding/json`（v1）のままなら意味論は不変（大小区別しない・重複後勝ち）。
>   `encoding/json/v2` を自分で import した時だけ厳格化。1.28 でも書き換え強制はなし。→ `knowledge/02`
> - v1 のままでも **Unmarshal は約1.8倍速**（v2バックエンドの恩恵。コード無改変）。
> - **gRPC は protobuf を使い `encoding/json` を通らない**ので JSON 変更は基本無関係。
> - 実質の唯一の注意は **最小 macOS が 13 Ventura**（CI/開発機/配布先）。→ `knowledge/00`
> - IDE の赤線は gopls を 1.27 でビルドし直せば解消（`GOTOOLCHAIN=go1.27.0 go install golang.org/x/tools/gopls@latest`）。

### 参考ソース

| - | URL |
| --- | --- |
| releaseNote | https://go.dev/doc/go1.27 |
| releaseBlog | https://go.dev/blog/go1.27 |
| JetBrains | https://blog.jetbrains.com/go/2026/08/20/ready-for-go-1-27-on-day-one/ |
| Linuxiac まとめ | https://linuxiac.com/go-1-27-released-with-generic-methods-json-v2-and-faster-memory-allocation/ |
| VictoriaMetrics tour | https://victoriametrics.com/blog/go-1-27/ |

### 特に気になっている部分について（勉強対象）

1. ジェネリックメソッド … `./genericMethods`
2. encoding/json/v2 … `./jsonV2`
3. goroutine リークプロファイル … `./goroutineLeak`
4. 構造体リテラルの柔軟化（埋め込み直接初期化）… `./structLiteral`
5. モダナイズ（`go fix`）と最新の書き方 … `./modernize` / `knowledge/05-modernize.md`, `knowledge/06-modernize-examples.md`
6. uuid パッケージ（新規）… `./uuid` / `knowledge/07-uuid.md`
7. Minor changes（strings.CutLast / url.Clone / rand.N / synctest.Sleep）… `./minorChanges` / `knowledge/08-minor-changes.md`

### 機能面以外のメモ（メモリ / 処理速度 / その他）

言語機能以外で押さえておきたい変更点。数値・フラグ名つき。

#### 🧠 メモリ

- **サイズ特化アロケータ（size-specialized malloc）**
  - コンパイラが割り当てサイズごとに専用ルーチンを生成。
  - **80バイト未満の小さい割り当てを最大30%高速化**。割り当てが多いプログラムで**全体約1%改善**。
  - 代償は**バイナリ +約60KB**（ワークロード非依存）。
  - 無効化: `GOEXPERIMENT=nosizespecializedmalloc`（1.28で削除予定＝実質恒久化）。

#### ⚡ 処理速度 / GC

- **Green Tea GC が恒久化**
  - GC自体は1.26でデフォルト化済み（小オブジェクトのマーク/スキャンを局所性とCPUスケーラビリティで改善）。
  - **GCオーバーヘッドを実プログラムで10〜40%削減**。Ice Lake / Zen4 以降のベクトル実装でさらに約10%。
  - **1.27でオプトアウト `GOEXPERIMENT=nogreenteagc` が削除**され、無効化できなくなった。
- **encoding/json（v2バックエンド）**: Marshalは従来同等、**Unmarshalが大幅高速化**。無効化: `GOEXPERIMENT=nojsonv2`。
- **compress/flate**: 圧縮が高速化。ただし**出力バイト列が1.26と変わりうる**
  （`archive/zip` / `compress/gzip` / `compress/zlib` / `image/png` の出力にも波及）。→ ゴールデンテスト注意。

#### 🔧 コンパイラ / リンカ

- **クロージャ命名の簡素化**: 関数リテラルの生成名を簡潔化し、インライン化時も同名に。
  同一クロージャの複数インスタンスがコード共有できるようになりバイナリ縮小に寄与。※シンボル名依存テストは要更新。
- **`//line` の相対パス解決変更**: 相対ファイル名をディレクティブのあるファイルのディレクトリ基準で解決（`go/scanner` と一致）。絶対パスは影響なし。
- **リンカ macOS フラグ**: `-macos` / `-macsdk` 追加（`LC_BUILD_VERSION` のOS/SDK指定）。デフォルトは最古 macOS 13.0.0 / SDK 26.2.0。

#### 👀 ランタイム / 観測性

- **goroutineleak プロファイル正式化**: `runtime/pprof`・`/debug/pprof/goroutineleak`（→ `knowledge/03-goroutine-leak.md`）。
- **traceback に pprof ラベル**: 1.27以降のモジュールで、panic 等の traceback ヘッダ行に `runtime/pprof` の goroutine ラベルが載る → 障害解析が楽に。無効化 `GODEBUG=tracebacklabels=0`（**恒久オプション**）。
- **`asynctimerchan` GODEBUG 削除**: `time` 由来チャネルは常に unbuffered（同期）。1.23以降の既定挙動が確定。

#### 📦 その他 / 互換性

- **最小 macOS が 13 Ventura** に（ビルド / CI / 配布ターゲット注意）。
- **bzr（Bazaar）サポート廃止**。
- 削除された GODEBUG 群: `tlsunsafeekm` / `tlsrsakex` / `tls3des` / `tls10server` / `x509keypairleaf` / `gotypesalias` など。

> 出典: [Go 1.27 Release Notes](https://go.dev/doc/go1.27) / [Go 1.26 Release Notes（Green Tea GC）](https://go.dev/doc/go1.26) / [InfoWorld: Green Tea GC](https://www.infoworld.com/article/4131097/go-1-26-unleashes-performance-boosting-green-tea-gc.html)

### 実行

```bash
cd go1.27.0
go run ./genericMethods
go run ./jsonV2
go run ./structLiteral
go run ./goroutineLeak
go run ./modernize
go run ./uuid
go run ./minorChanges
go test ./minorChanges -v    # synctest.Sleep の例
go test ./goroutineLeak -v   # リーク検査(CI)テンプレ
go fix -diff ./...           # モダナイズ差分のドライラン
```

> ツールチェーンは `GOTOOLCHAIN=auto`。このディレクトリは go.mod の `go 1.27.0` に従い
> 自動で 1.27 を取得・使用する（他のバージョンディレクトリや別プロジェクトには影響しない）。
