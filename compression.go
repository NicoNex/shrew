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
	"strings"

	"github.com/mholt/archiver/v3"
)

type tar interface {
	Archive(src []string, dest string) error
}

type Compression struct {
	ext string
	arc tar
}

func extract(src string, dest string) error {
	return archiver.Unarchive(src, dest)
}

func NewCompression(name string) Compression {
	var t tar
	var e string

	switch name {
	case "targz":
		e = ".tar.gz"
		tmp := archiver.NewTarGz()
		tmp.OverwriteExisting = true
		t = tmp
	case "tarzstd":
		e = ".tar.zst"
		tmp := archiver.NewTarZstd()
		tmp.OverwriteExisting = true
		t = tmp
	default:
		e = ".zip"
		tmp := archiver.NewZip()
		tmp.OverwriteExisting = true
		t = tmp
	}

	return Compression{
		ext: e,
		arc: t,
	}
}

func (c Compression) GetFilename(fname string) string {
	return strings.Join([]string{fname, c.ext}, "")
}

func (c Compression) Compress(src string, dest string) error {
	return c.arc.Archive([]string{src}, dest)
}
