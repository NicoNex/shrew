package main

import (
	"strings"

	"github.com/mholt/archiver/v3"
)

type archCreator interface {
	Archive(src []string, dest string) error
}

type Compression struct {
	ext string
	arc archCreator
}

func extract(src string, dest string) error {
	return archiver.Unarchive(src, dest)
}

func NewCompression(name string) Compression {
	switch name {
	case "targz":
		return Compression{
			ext: ".tar.gz",
			arc: archiver.NewTarGz(),
		}
	case "tarzstd":
		return Compression{
			ext: ".tar.zst",
			arc: archiver.NewTarZstd(),
		}
	default:
		return Compression{
			ext: ".zip",
			arc: archiver.NewZip(),
		}
	}
}

func (c Compression) GetFilename(fname string) string {
	return strings.Join([]string{fname, c.ext}, "")
}

func (c Compression) Compress(src string, dest string) error {
	return c.arc.Archive([]string{src}, dest)
}
