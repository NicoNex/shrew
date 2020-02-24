package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

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
