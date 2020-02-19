package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
)

var cfg Config
var wg sync.WaitGroup

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

// TODO: refactor this function.
func saveFile(archive string, fname string, data []byte, c chan Item) {
	defer wg.Done()
	path := fmt.Sprintf("%s/%s", cfg.Path, archive)
	if _, err := os.Stat(path); err != nil {
		err := os.MkdirAll(path, 0644)
		if err != nil {
			c <- Item{
				fname,
				archive,
				NewStatus(err),
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
			NewStatus(err),
		}
		log.Println(err)
		return
	}

	c <- Item{
		fname,
		archive,
		NewStatus(nil),
	}
}

func getDirEntries(dirpath string) ([]string, error) {
	var ret []string

	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return []string{}, err
	}
	for _, f := range files {
		ret = append(ret, f.Name())
	}

	return ret, nil
}

func fetchArchives(path string) ([]Archive, error) {
	var archives []Archive

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return []Archive{}, err
	}

	for _, f := range files {
		if f.IsDir() {
			var tmp Archive

			name := f.Name()
			tmpname := fmt.Sprintf("%s/%s", path, name)
			fnames, err := getDirEntries(tmpname)
			if err != nil {
				tmp = NewArchiveErr(name, err)
			} else {
				tmp = NewArchive(name, fnames)
			}
			archives = append(archives, tmp)
		}
	}

	return archives, nil
}

func showHomePage(w http.ResponseWriter, r *http.Request) {
	var response string

	archives, err := fetchArchives(cfg.Path)
	if err != nil {
		response = GetErrResponse(err)
	} else {
		response = GetHomeResponse(archives)
	}
	fmt.Fprintf(w, response)
}

func getItems(itemch chan Item, outch chan []Item) {
	var items []Item

	for i := range itemch {
		items = append(items, i)
	}
	outch <- items
}

// TODO: refactor this function.
func handleUpload(w http.ResponseWriter, r *http.Request) {
	var err error
	var name string
	var response string
	var ichan = make(chan Item)
	var outchan = make(chan []Item, 1)

	if r.Method != "POST" {
		err := errors.New("Invalid request")
		response = GetErrResponse(err)
		goto write_response
	}

	// 1Mb in memory the rest on disk.
	r.ParseMultipartForm(1048576)
	name, err = getQuery("name", r.URL.RawQuery)
	if err != nil {
		response = GetErrResponse(err)
		goto write_response
	}

	if r.MultipartForm == nil {
		err := errors.New("No file provided")
		response = GetErrResponse(err)
		goto write_response
	}

	go getItems(ichan, outchan)
	for _, headers := range r.MultipartForm.File {
		for _, h := range headers {
			tmp, err := h.Open()
			if err != nil {
				log.Println(err)
				continue
			}

			filename := h.Filename
			filedata, err := ioutil.ReadAll(tmp)
			tmp.Close()
			if err != nil {
				log.Println(err)
				continue
			}
			wg.Add(1)
			go saveFile(name, filename, filedata, ichan)
		}
	}
	wg.Wait()
	close(ichan)
	response = GetUploadResponse(<-outchan)
	close(outchan)
write_response:
	fmt.Fprint(w, response)
}

// TODO: complete this.
func handleDownload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Function coming soon...")
}

func getConfig() Config {
	var cfgpath string

	if runtime.GOOS == "windows" {
		home := os.Getenv("UserProfile")
		cfgpath = fmt.Sprintf("%s/.shrew/config.toml", home)
	} else {
		home := os.Getenv("HOME")
		cfgpath = fmt.Sprintf("%s/.config/shrew/config.toml", home)
	}

	cfg, err := loadConfig(cfgpath)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func main() {
	cfg = getConfig()
	http.HandleFunc("/", showHomePage)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/download", handleDownload)

	port := fmt.Sprintf(":%d", cfg.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}
