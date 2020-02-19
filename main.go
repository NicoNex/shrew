package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func saveFile(archive string, fname string, data []byte, c chan Item) {
	defer wg.Done()
	path := filepath.Join(cfg.Path, archive)
	if _, err := os.Stat(path); err != nil {
		if err := os.MkdirAll(path, 0644); err != nil {
			c <- NewItem(fname, archive, err)
			log.Println(err)
			return
		}
	}

	path = filepath.Join(path, fname)
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		c <- NewItem(fname, archive, err)
		log.Println(err)
		return
	}

	c <- NewItem(fname, archive, nil)
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
			tmpname := filepath.Join(path, name)
			fmt.Println(tmpname)
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

func collectItems(itemch chan Item, outch chan []Item) {
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

	go collectItems(ichan, outchan)
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

// TODO: complete this.
func handleDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Function coming soon...")
}

func getConfig() Config {
	var cfgpath string

	if runtime.GOOS == "windows" {
		cfgpath = filepath.Join(os.Getenv("UserProfile"), ".shrew/config.toml")
	} else {
		cfgpath = filepath.Join(os.Getenv("HOME"), ".config/shrew/config.toml")
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
	http.HandleFunc("/delete", handleDelete)

	port := fmt.Sprintf(":%d", cfg.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}
