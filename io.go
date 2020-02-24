package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Returns true if a file or directory exists.
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Saves a file on the disk and sends the info in the channel provided.
func saveFile(archive string, fname string, data []byte, c chan Item) {
	defer wg.Done()
	path := filepath.Join(cfg.Path, archive)
	if !exists(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			c <- NewItem(fname, archive, "", err)
			log.Println(err)
			return
		}
	}

	path = filepath.Join(path, fname)
	err := ioutil.WriteFile(path, data, 0755)
	if err != nil {
		c <- NewItem(fname, archive, "", err)
		log.Println(err)
		return
	}

	c <- NewItem(fname, archive, path, nil)
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

// Fetches all the archives from disk and returns them.
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
