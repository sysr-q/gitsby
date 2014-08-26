package main

import (
	"encoding/json"
	"io/ioutil"
)

type configRepo struct {
	Url       string `json:"url"`
	Directory string `json:"directory"`
	Hidden    bool   `json:"hidden"`
	Type      string `json:"type"`
	Command   CommandInfo `json:"command"`
}

type Config struct {
	Repos map[string]Repo
}

func ParseConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := Config{Repos: make(map[string]Repo)}

	// Temporary variable so we can unmarshal correctly.
	var c struct{
		Repos []configRepo `json:"repos"`
	}

	if err = json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	for _, repo := range c.Repos {
		// TODO: switch repo.Type -> not just Git support.
		r := Git{
			Url: repo.Url,
			directory: repo.Directory,
			Hidden: repo.Hidden,
			Command: repo.Command,
		}

		owner, name := r.Name()
		config.Repos[owner + "/" + name] = r
	}

	return &config, nil
}
