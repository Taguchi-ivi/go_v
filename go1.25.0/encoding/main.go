package main

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// 現在experimentalなのでGOEXPERIMENT=jsonv2が必要になります。
// これは結構でかい変更な気がする
// tagが追加されている。これによって存在しないものが送られた場合にlogに出力することができるようになった。
// 変な値を検知できるかも?
// $ GOEXPERIMENT=jsonv2 go run main.go

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	// Address   string         `json:",inline"` // インラインで表示される
	Age       int            `json:"age"`
	BirthDate time.Time      `json:"birth_date,format:'2006-01-02'"` // 日付のみのフォーマット
	Password  string         `json:"password"`
	Extra     map[string]any `json:",unknown"` // 不明なフィールド受け取り
}

func main() {
	// これはslogの勉強
	// slog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.Info("Hello, world!", "key1", "value1", "key2", 42, "key3", true)
	slog.InfoContext(context.Background(), "Contextual log", "contextKey", "contextValue")

	// ここからencodingの勉強
	src := []byte(`{
		"name": "Alice",
		"birth_date": "2025-07-15",
		"street": "123 Main St",
		"city": "Wonderland",
		"hobby": "Adventuring",
		"friends": [{"name": "Bob"}, {"name": "Cindy"}]
 	}`)

	var user User
	if err := json.Unmarshal(src, &user); err != nil {
		slog.Error("Failed to unmarshal JSON", "error", err)
		os.Exit(1)
	}
	slog.Info("Unmarshaled User", "user", user)
	slog.Info("存在しないフィールドについて", "unknownField", user.Extra)

	// outputも
	out, err := json.Marshal(user, jsontext.WithIndent("    "))
	if err != nil {
		slog.Error("Failed to marshal User", "error", err)
		os.Exit(1)
	}
	// slog.Info("Marshaled User", "output", string(out))
	fmt.Println(string(out))
}
