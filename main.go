/*
 * Shrew - a little shrew that stores files or backups with a REST API.
 * Copyright (C) 2020  Nicolò Santamaria
 * 
 * Shrew is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * 
 * Shrew is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU General Public License
 * along with shrew.  If not, see <https://www.gnu.org/licenses/>.
 */

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

// Returns the value of an url raw query or error if missing.
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

// Serves the main data in the body of the request.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	var response string

	archives, err := fetchArchives(cfg.Path)
	if err != nil {
		response = GetStatusResponse(err)
	} else {
		response = GetRootResponse(archives)
	}
	fmt.Fprintln(w, response)
}

// Collects all the Item objects sent in the itemch and sends the
// reulting array in the outch.
func collectItems(itemch chan Item, outch chan []Item) {
	var items []Item

	for i := range itemch {
		items = append(items, i)
	}
	outch <- items
}

// Handles the upload of a file or archive and stores it.
func putHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var name string
	var response string
	var ichan = make(chan Item)
	var outchan = make(chan []Item, 1)

	if r.Method != "POST" {
		err := errors.New("Invalid request")
		response = GetStatusResponse(err)
		goto write_response
	}

	// 1Mb in memory the rest on disk.
	r.ParseMultipartForm(1048576)
	name, err = getQuery("archive", r.URL.RawQuery)
	if err != nil {
		response = GetStatusResponse(err)
		goto write_response
	}

	if r.MultipartForm == nil {
		err := errors.New("No file provided")
		response = GetStatusResponse(err)
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
	response = GetItemsResponse(<-outchan)
	close(outchan)

write_response:
	fmt.Fprintln(w, response)
}

// TODO: refactor this.
func getHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var comp string
	var archive string
	var response string
	var archpath string
	var archfile string
	var cachepath string
	var cmp Compression

	archive, err = getQuery("archive", r.URL.RawQuery)
	if err != nil {
		log.Println(err)
		response = GetStatusResponse(err)
		goto write_response
	}

	cachepath = filepath.Join(cfg.Path, ".cache")
	if !exists(cachepath) {
		if err = os.MkdirAll(cachepath, 0755); err != nil {
			log.Println(err)
			response = GetStatusResponse(err)
			goto write_response
		}
	}

	comp, err = getQuery("compression", r.URL.RawQuery)
	if err != nil {
		comp = cfg.Compression
	}

	archpath = filepath.Join(cfg.Path, archive)
	if !exists(archpath) {
		err = fmt.Errorf("%s: archive not found", archive)
		log.Println(err)
		response = GetStatusResponse(err)
		goto write_response
	}

	cmp = NewCompression(comp)
	archfile = cmp.GetFilename(archive)
	cachepath = filepath.Join(cachepath, archfile)
	if !exists(cachepath) || !isUpToDate(archpath, cachepath) {
		err = cmp.Compress(archpath, cachepath)
		if err != nil {
			log.Println(err)
			response = GetStatusResponse(err)
			goto write_response
		}
	}

	http.ServeFile(w, r, cachepath)
	// response = GetItemsResponse([]Item{
	// 	Item{
	// 		Name:    archfile,
	// 		Archive: archive,
	// 		Hash:    checksum(cachepath),
	// 		Status:  NewStatus(nil),
	// 	},
	// })

write_response:
	fmt.Fprintln(w, response)
}

// Handles the removal of a file or archive.
func delHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var fnames string
	var archive string
	var response string
	var archpath string

	archive, err = getQuery("archive", r.URL.RawQuery)
	if err != nil {
		response = GetStatusResponse(err)
		goto write_response
	}

	archpath = filepath.Join(cfg.Path, archive)
	fnames, err = getQuery("files", r.URL.RawQuery)
	if err != nil {
		err := os.RemoveAll(archpath)
		response = GetStatusResponse(err)
	} else {
		namearr := strings.Split(fnames, ",")
		inch := make(chan Item, len(namearr))
		outch := make(chan []Item, 1)

		go collectItems(inch, outch)
		for _, n := range namearr {
			wg.Add(1)
			go func(name string, ch chan Item) {
				defer wg.Done()
				path := filepath.Join(archpath, name)
				err := os.Remove(path)
				ch <- NewItem(name, archive, path, err)
			}(n, inch)
		}
		wg.Wait()
		close(inch)
		response = GetItemsResponse(<-outch)
		close(outch)
	}

write_response:
	fmt.Fprintln(w, response)
}

// Returns a Config object with values from the config file.
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
	var msg = `
  __QQ
 (_)_">
_)      shrew running..

`
	fmt.Print(msg)
	cfg = getConfig()

	if !exists(cfg.Path) {
		log.Println("Archive directory not found, creating it.")
		err := os.MkdirAll(cfg.Path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/put", putHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/del", delHandler)

	port := fmt.Sprintf(":%d", cfg.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}
