# Minor changes to the library（実務で使いそうな抜粋）

Go 1.27 の細かい標準ライブラリ変更のうち、実務で当たりそうなものだけ。
すべて 1.27 の `go doc` で実在を確認済み。コードは `../minorChanges/`。

## strings.CutLast / bytes.CutLast

```go
func CutLast(s, sep string) (before, after string, found bool)
```

既存 `Cut` は「最初の」区切りで分割。`CutLast` は**「最後の」区切り**で分割。

```go
strings.CutLast("src/pkg/util/helper.go", "/") // "src/pkg/util", "helper.go", true
strings.CutLast("archive.tar.gz", ".")         // "archive.tar", "gz", true
```

パスとファイル名、名前と拡張子など「末尾側で切りたい」場面にちょうどよい。

## net/url URL.Clone / Values.Clone

```go
func (u *URL) Clone() *URL
func (vs Values) Clone() Values
```

ディープコピー。HTTP ミドルウェアで**元の URL を壊さずにクエリを足す**等で便利。

```go
req := base.Clone()
q := req.Query(); q.Set("page", "2"); req.RawQuery = q.Encode()
// base は無傷
```

## math/rand/v2 Rand.N（ジェネリックメソッド）

```go
func (r *Rand) N[Int intType](n Int) Int
```

**1.27 のジェネリックメソッド機能を標準ライブラリ自身が使っている実例**（→ [[01-generic-methods]]）。
任意の整数型で `[0, n)` の乱数を、型変換なしで得られる。

```go
r.N(int32(100))       // int32 で返る
r.N(uint(1_000_000))  // uint で返る
```

## testing/synctest Sleep

```go
func Sleep(d time.Duration)
```

synctest バブル内の**仮想時計を d だけ進め**、（自分以外の）全 goroutine が
durably-blocked になるまで待つ。`time.Sleep + synctest.Wait` を 1 行にしたもの。
**実時間を待たずに「N時間後」の挙動をテスト**できる。

```go
synctest.Test(t, func(t *testing.T) {
    go func() { time.Sleep(2 * time.Hour); close(done) }()
    synctest.Sleep(1 * time.Hour) // まだ未完了
    synctest.Sleep(1 * time.Hour) // 合計2時間 -> 完了
})
```

コード: `../minorChanges/main_test.go`（`go test ./minorChanges -v`）

## その他（デモ無し・知っておくと良い）

| パッケージ | 変更 | 実務での意味 |
| --- | --- | --- |
| net/http | `Server.MaxHeaderValueCount` / `DefaultMaxHeaderValueCount` | ヘッダ数の上限でDoS耐性を調整 |
| net/http | HTTP/1 の `Response.Body` を close 時に自動ドレイン | 接続再利用が効きやすくなる |
| net/http | `Server.DisableClientPriority`（HTTP/2 優先度） | 挙動調整ノブ |
| database/sql | `ConvertAssign` 関数、`driver.RowsColumnScanner` | ドライバ実装者向け |
| crypto/x509 | `SystemCertPool` が `SSL_CERT_FILE` / `SSL_CERT_DIR` を尊重 | コンテナで証明書パス指定が効く |
| net | `UnixConn` の read が `io.EOF` を直接返す | Unixソケットの EOF 判定が素直に |
| unicode | Unicode 15 → 17 | 新しい文字・絵文字に対応 |

## 関連
- コード: `../minorChanges/`
- [[00-overview]] / [[01-generic-methods]]
