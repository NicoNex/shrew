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
