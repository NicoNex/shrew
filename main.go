package main

import (
	"os"
	"fmt"
	"log"
	"sync"
	"errors"
	"runtime"
	"strings"
	"net/http"
	"io/ioutil"
)

var cfg Config
var wg sync.WaitGroup

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

func saveFile(archive string, fname string, data []byte, c chan Item) {
	defer wg.Done()
	var result Item

	path := fmt.Sprintf("%s/%s", cfg.Path, archive)
	if _, err := os.Stat(path); err != nil {
		err := os.MkdirAll(path, 0644)
		if err != nil {
			c <- Item{
				fname,
				archive,
				false,
				err.Error(),
			}
			log.Println(err)
			return
		}
	}

	path = fmt.Sprintf("%s/%s", path, fname)
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		c <- Item{
			fname,
			archive,
			false,
			err.Error(),
		}
		log.Println(err)
		return
	}

	result = Item{
		Name: fname,
		Archive: archive,
		Ok: true,
	}

	c <- result
}

// TODO: make it fetch the items from somewhere.
func showHomePage(w http.ResponseWriter, r *http.Request) {
	res := GetItemsResponse([]Item{})
	fmt.Fprintf(w, res)
}

// TODO: maybe refactor this function.
func handleUpload(w http.ResponseWriter, r *http.Request) {
	var err error
	var name string
	var response string
	var items []Item
	var fileCount int
	var ichan = make(chan Item)
	var outchan = make(chan []Item, 1)

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

	// TODO: put this into a separate non-anonymous function.
	go func(ichan chan Item, outchan chan []Item) {
		var items []Item
		
		for i := range ichan {
			items = append(items, i)
		}
		outchan <- items
	}(ichan, outchan)

	for _, headers := range r.MultipartForm.File {
		for _, h := range headers {
			tmp, err := h.Open()
			if err != nil {
				log.Println(err)
				continue
			}
			defer tmp.Close()

			filename := h.Filename
			filedata, err := ioutil.ReadAll(tmp)
			if err != nil {
				log.Println(err)
				continue
			}
			fileCount++
			wg.Add(1)
			go saveFile(name, filename, filedata, ichan)
		}
	}
	wg.Wait()
	close(ichan)
	items = <-outchan
	close(outchan)

	response = GetItemsResponse(items)
write_response:
	fmt.Fprint(w, response)
}

// TODO: complete this.
func handleDownload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Function coming soon...")
}

func main() {
	var cfgpath string

	if runtime.GOOS == "windows" {
    	home := os.Getenv("UserProfile")
    	cfgpath = fmt.Sprintf("%s/.shrew/config.toml", home)
    } else {
    	home := os.Getenv("HOME")
    	cfgpath = fmt.Sprintf("%s/.config/shrew/config.toml", home)
    }

    {
    	var err error
	    cfg, err = loadConfig(cfgpath)
	    if err != nil {
	    	log.Fatal(err)
	    }
	}

    http.HandleFunc("/", showHomePage)
    http.HandleFunc("/upload", handleUpload)
    http.HandleFunc("/download", handleDownload)

    port := fmt.Sprintf(":%d", cfg.Port)
    log.Fatal(http.ListenAndServe(port, nil))
}
