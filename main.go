package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
)

func check(e error) bool {
	var isErr = (e != nil)
	if isErr {
		fmt.Println(e)
	}
	return isErr
}

func getQuery(name string, rawQuery string) (string, error) {
	var queries = strings.Split(rawQuery, "&")
	for _, q := range queries {
		tokens := strings.Split(q, "=")
		if tokens[0] == name {
			return tokens[1], nil
		}
	}
	return "", fmt.Errorf("%s: query not found", name)
}

// TODO: make it fetch the items from somewhere.
func showHomePage(w http.ResponseWriter, r *http.Request) {
	res := GetHomeResponse([]Item{})
	fmt.Fprintf(w, res)
}

// TODO: maybe refactor this function.
func handleUpload(w http.ResponseWriter, r *http.Request) {
	var err error
	var name string
	var version string
	var response string

	name, err = getQuery("name", r.URL.RawQuery)
	if check(err) {
		response = GetErrResponse(err)
		goto write_response		
	}

	version, err = getQuery("version", r.URL.RawQuery)
	if check(err) {
		response = GetErrResponse(err)
		goto write_response
	}

	response = GetUploadResponse(name, version, "put path here", err)
write_response:
	fmt.Fprint(w, response)
}

func main() {
    http.HandleFunc("/", showHomePage)
    http.HandleFunc("/upload", handleUpload)

    log.Fatal(http.ListenAndServe(":8081", nil))
}
