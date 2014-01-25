package server

import (
	"fmt"
	"errors"
	"io/ioutil"
	"encoding/json"
	"github.com/hoisie/web"
	"github.com/plausibility/gitsby/git"
	"github.com/plausibility/gitsby/util"
)

var Server *web.Server
var Config GitsbyConfig

func readConfig(path string) (GitsbyConfig, error) {
	var c GitsbyConfig
	if exists, _ := util.FileExists(path); !exists {
		return c, errors.New(fmt.Sprintf("Config file %s doesn't exist.", path))
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		// programming
		return c, err
	}
	if jerr := json.Unmarshal(data, &c); jerr != nil {
		// just try and stop me
		return c, jerr
	}
	// Set up a map of name -> repo.
	c.ReposActive = make(map[string]git.Repo)
	for _, repo := range c.Repos {
		owner, name := repo.Name()
		c.ReposActive[owner + "/" + name] = repo
	}
	return c, nil
}

func Setup(config *string) error {
	// Load our config file.
	if conf, loadErr := readConfig(*config); loadErr != nil {
		return loadErr
	} else {
		Config = conf
	}
	// Set up our web.go server
	Server = web.NewServer()
	return nil
}
