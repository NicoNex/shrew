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
