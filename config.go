package main

import "github.com/BurntSushi/toml"

type Config struct {
	Path        string `toml:"archive_path"`
	Port        int    `toml:"connection_port"`
	Overwrite   bool   `toml:"allow_overwrite"`
	SaveCache   bool   `toml:"cache_archive"`
	Compression string `toml:"compression"`
}

func loadConfig(path string) (Config, error) {
	var cfg Config

	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
