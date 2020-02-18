package main

import (
	"fmt"
	"log"
	"errors"
	"strings"
	"net/http"
	"io/ioutil"
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

	if r.Method != "POST" {
		err = errors.New("Invalid request")
		response = GetErrResponse(err)
		goto write_response
	}

	// 1Mb in memory the rest on the disk.
	r.ParseMultipartForm(1048576)

	name, err = getQuery("name", r.URL.RawQuery)
	if check(err) {
		response = GetErrResponse(err)
		goto write_response		
	}

	if r.MultipartForm == nil {
		err = errors.New("No file provided")
		response = GetErrResponse(err)
		goto write_response
	}

	for _, headers := range r.MultipartForm.File {
		fmt.Println("sas mike")
		for _, h := range headers {
			tmp, err := h.Open()
			if err != nil {
				response = GetErrResponse(err)
				goto write_response
			}
			defer tmp.Close()

			fmt.Println(h.Filename)
			filename := h.Filename
			filedata, err := ioutil.ReadAll(tmp)
			if check(err) {
				response = GetErrResponse(err)
				goto write_response
			}
			fmt.Println(filename, string(filedata))
			// TODO: add the logic to save the file where specified.
		}
	}

	response = GetUploadResponse(name, version, "put path here", err)
write_response:
	fmt.Fprint(w, response)
}

// TODO: complete this.
func handleDownload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Function coming soon...")
}

// TODO: load the config file.
func main() {
    http.HandleFunc("/", showHomePage)
    http.HandleFunc("/upload", handleUpload)
    http.HandleFunc("/download", handleDownload)

    log.Fatal(http.ListenAndServe(":8081", nil))
}
