package main

import (
	"log"
	"net/http"
)

func main() {
	// CSRFの対応。echoとかだとCORSWithConfigで設定できるのかな??

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	cop := http.NewCrossOriginProtection()
	cop.AddTrustedOrigin("https://example.com")
	cop.AddTrustedOrigin("https://*.example.com")
	cop.AddTrustedOrigin("localhost:8080")

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", cop.Handler(mux))

}
