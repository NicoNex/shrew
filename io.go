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
	"crypto/sha256"
	"io"
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

func checksum(fpath string) string {
	file, err := os.Open(fpath)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}

	return string(hash.Sum(nil))
}

func isUpToDate(arname string, ccname string) bool {
	archinfo, err := os.Stat(arname)
	if err != nil {
		log.Println(err)
		return false
	}

	cacheinfo, err := os.Stat(ccname)
	if err != nil {
		log.Println(err)
		return false
	}

	cachetime := cacheinfo.ModTime()
	archtime := archinfo.ModTime()
	return cachetime.After(archtime)
}
