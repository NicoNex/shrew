package main

import (
	"fmt"
	"log"
	"html"
	"strings"
	"net/http"
)

func showHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	var queries = strings.Split(r.URL.RawQuery, "&")
	
	fmt.Fprintf(w, "Upload test: %v", r.URL.RawQuery)
}

func main() {
    http.HandleFunc("/", showHomePage)
    http.HandleFunc("/upload", handleUpload)

    log.Fatal(http.ListenAndServe(":8081", nil))
}
