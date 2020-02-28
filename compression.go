package main

import "github.com/mholt/archiver/v3"

type Format int8
const (
	Zip 	 Format = iota
	TarGz
	TarZstd
)

type Archiver interface {
	Archive(src []string, dest string) error
	Unarchive(src, dest string) error
}

func compress(src string, dest string, f Format) error {
	var arc Archiver

	switch f {
	case Zip:
		arc = archiver.NewZip()
	case TarZstd:
		arc = archiver.NewTarZstd()
	case TarGz:
		arc = archiver.NewTarGz()
	}

	return arc.Archive([]string{src}, dest)
}

func extract(src string, dest string) error {
	return archiver.Unarchive(src, dest)
}
