# goroutine リークプロファイル（goroutineleak）

## 概要

`runtime/pprof` に `"goroutineleak"` プロファイルが正式追加。
チャネル / `sync.Mutex` / `sync.Cond` などの同期プリミティブで
**永久にブロックし、二度と復帰しない goroutine** を検出する。

- Uber の Vlad Saioc の貢献。実験段階（GODEBUG）を経て 1.27 で正式化。
- 完全検出は原理的に不可能だが、大きなクラスのリークを拾える。

## 判定の核心：GC の「到達可能性」（誤検知しない仕組み）

**単にブロックしているだけではリーク扱いされない。** ブロック対象（チャネル/Mutex等）が
「もうどこからも到達できない＝誰も二度と signal できない」状態のときだけリークと判定する。
メモリの GC と同じ発想を goroutine に適用したもの。

実機で確認（`../goroutineLeak/main.go`）: 3つブロック中でもリーク報告は2件だけ。

| goroutine | 待ち対象 | 判定 |
| --- | --- | --- |
| `realLeak` ×2 | ローカルの `ch`（他から到達不能） | **リーク** |
| `blockedButReachable` | グローバルの `reachable`（到達可能） | リークでない（誤検知しない） |

```
goroutineleak profile: total 2      // 3ブロック中でも 2 だけ
#  main.realLeak+0x2f  .../goroutineLeak/main.go:37
```

→ 「まだ signal されうる正常な待ち」を巻き込まないので、**プロダクションでも安心して読める**。

## 使い方

```go
p := pprof.Lookup("goroutineleak")
p.WriteTo(os.Stdout, 1) // debug=1 でスタック付き
```

- `net/http/pprof` を import していれば `/debug/pprof/goroutineleak` でも取得可能。
- **`WriteTo` がリーク検出のために GC サイクルを走らせる。**
  そのため `Count()` を WriteTo 前に読むと 0 のことがある（検出は書き出し時に走る）。

## 実測（`go run ./goroutineLeak`）

送信されないチャネルを待つ goroutine を 2 本作った場合:

```
--- goroutineleak profile (debug=1) ---
goroutineleak profile: total 2
2 @ 0x... 0x... 0x... 0x... 0x...
#	0x...	main.leak+0x2f	.../goroutineLeak/main.go:27
```

リーク発生箇所（`main.leak`）とスタックが出るので原因特定に使える。

## 実務での回し方（テストは必須ではない）

| 方式 | テスト要否 | 使う場面 |
| --- | --- | --- |
| **C. 本番サーバで観測** | **不要** | 稼働中サーバが今リークしてるか見る |
| **A. CI で明示アサート** | 少し要（ヘルパー1個） | マージ前にリーク混入を止める（予防の本命） |
| **B. synctest で自動検知** | synctest を使うだけ | 時間依存コードのテストで、リークが自動で落ちる |

**C. 本番（テスト不要）**: `net/http/pprof` を import してエンドポイントを生やすだけ。
```go
import _ "net/http/pprof"
go http.ListenAndServe("localhost:6060", nil)
// 稼働中: curl http://localhost:6060/debug/pprof/goroutineleak?debug=1
```

**A. CI アサート**: プロファイルを走らせて件数で落とす（`../goroutineLeak/leakcheck_test.go`）。
```go
p := pprof.Lookup("goroutineleak")
p.WriteTo(io.Discard, 1)          // 検出(GC)をトリガ。Count はこの後に有効（実測: 2 を返す）
if n := p.Count(); n > 0 { t.Fatalf("leaked %d", n) }
```

**B. synctest 自動検知**: `synctest.Test` のバブル終了時にブロック残があると
`panic: deadlock: main bubble goroutine has exited but blocked goroutines remain` で失敗。
→ **専用アサートを書かなくても**リークが落ちる（実測で確認）。

> おすすめ: 予防は **A/B を CI**、実地調査は **C を本番** に、の二段構え。

## 実務での使いどころ

- 長時間稼働サーバのリークハント（テストの `testing/synctest` と併用すると再現しやすい → [[00-overview]]）。
- CI で「リーク 0 であること」を assert する回帰テストに組み込む。

## 関連

- コード: `../goroutineLeak/main.go`
