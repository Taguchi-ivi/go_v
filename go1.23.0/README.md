## v.1.23.0 の勉強

ちょっとリポジトリの切り方ミスった。
今後全ての v を一つのリポジトリで管理できるようにしたい。まあ、今後

## 資料

- [Release Notes](https://tip.golang.org/doc/go1.23)
- [blog](https://go.dev/blog/go1.23)
- [ユーチャーさんの記事(slices/map)](https://future-architect.github.io/articles/20240719a/)
- [ユーチャーさんの記事(range over func,iter)](https://future-architect.github.io/articles/20240718a/)

## playground

https://go.dev/play/

## メモと追記

- go command
  go mod tidy diff が使えるようになった。
  これによって、ファイルを修正せず、代わりに必要な変更を統一された diff として出力する
  更新が必要な場合は 0 以外のコードで終了する

- godebug
  便利そうだけど用途まだ不明 そもそも godebug フラグを使ってない。
  https://tip.golang.org/doc/godebug

- Timer, Ticker
  今まで stop が呼び出されないと gc(ガベージコレクション)の対象にならなかったが stop しなくても対象になる。とても良い。
  メモリの開放的な話だね。

- go linkname
  廃止にしたいが、サードパーティ系で使われていて、消すにも消せない。
  そもそもこれは、関数の小文字から始まる private の関数を public みたいにできるもの。
  使わないほうが良い。
