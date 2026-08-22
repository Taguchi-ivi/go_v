# Go 1.27 概要 & 移行チェックリスト

> リリース: 2026年8月 / このリポジトリでの検証ツールチェーン: `go1.27.0 darwin/arm64`

## 目玉

| カテゴリ | 変更 |
| --- | --- |
| 言語 | **ジェネリックメソッド**（メソッドが独自の型パラメータを持てる） |
| 言語 | **構造体リテラル**で埋め込みフィールドを直接初期化 |
| 言語 | 関数型推論の一般化（ジェネリック関数を関数型へ代入/変換する全文脈で推論が効く） |
| stdlib | **encoding/json/v2** + **encoding/json/jsontext**（新JSONエンジン、厳格デフォルト） |
| stdlib | ポスト量子暗号 `crypto/mldsa`（ML-DSA / FIPS 204）、TLS で ML-DSA / ML-KEM |
| stdlib | UUID パッケージ（新規） |
| runtime | 小サイズ割り当て（<80B）を最大30%高速化（サイズ特化ルーチン） |
| runtime | **goroutineleak** プロファイルが正式化（`runtime/pprof`） |
| tool | `go fix` にモダナイザ多数（`atomictypes` `embedlit` `slicesbackward` `unsafefuncs` 等） |
| tool | bzr（Bazaar）サポート廃止 |
| port | **最小 macOS が 13 Ventura** に引き上げ |

## GOEXPERIMENT オプトアウト（いずれも 1.28 で削除予定）

- `GOEXPERIMENT=nojsonv2` … JSON の内部実装を v1 に戻す
- `GOEXPERIMENT=nosizespecializedmalloc` … 割り当て高速化を無効化

## ツールチェーン運用メモ（このリポジトリでの実際）

- グローバル `GOTOOLCHAIN` を **`auto`** に設定済み。
  各ディレクトリが自分の `go.mod` の `go` 行に従ってツールチェーンを選ぶ。
  → `go1.27.0/` は 1.27、他は従来通り。**別プロジェクトは go.mod を上げるまで 1.27 に上がらない**。
- ローカルにツールチェーン未取得でも、`go.mod` が `go 1.27.0` なら
  `auto` 下で `go` 実行時に自動DLされる（`go: downloading go1.27.0`）。

## 別プロジェクトを 1.27 へ切り替える前のチェックリスト

| # | 確認 | 影響度 | 対応 |
| --- | --- | --- | --- |
| 1 | ビルド/CI/実行環境の macOS が **13 Ventura 以上**か | 高 | CIランナー・開発機・配布ターゲットを確認 |
| 2 | JSON のエラーメッセージ**文字列**に依存したテストがないか | 低 | v1 API のままなら意味論は不変（大小区別・重複後勝ちは維持）。変わるのはエラー文言程度。厳格化は `encoding/json/v2` へ**自分で移行した時だけ**（→ [[02-json-v2]]） |
| 3 | `compress/flate` の**出力バイト一致**や関数リテラルのシンボル名に依存するテスト | 低 | ゴールデンテストを更新 |
| 4 | 削除された GODEBUG（`asynctimerchan` など）に依存していないか | 低 | timer チャネルは unbuffered 前提（1.23 以降の既定挙動） |
| 5 | bzr でホストされた依存があるか | 低 | 事実上ほぼ無い |

> Go 1 互換保証があるので**基本はソース無改変でコンパイルできる**。上記に該当しなければ切替リスクは小さい。
