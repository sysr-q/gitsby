package server

import (
	"os"
	"fmt"
	"path"
	"errors"
	"os/user"
	"io/ioutil"
	"encoding/json"
	"github.com/hoisie/web"
)

type Repo struct {
	Url string `json:"url"`
	Directory string `json:"directory"`
}

type GitsbyConfig struct {
	Landing bool `json:"landing"`
	Repos []Repo `json:"repos"`
}

var Server *web.Server
var Config GitsbyConfig

func exists(path string) (bool, error) {
	    _, err := os.Stat(path)
	    if err == nil { return true, nil }
	    if os.IsNotExist(err) { return false, nil }
	    return false, err
}

func readConfig(path string) (GitsbyConfig, error) {
	var c GitsbyConfig
	if exists, _ := exists(path); !exists {
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
	if len(c.Repos) == 0 {
		// no really, just try stopping me
		return c, errors.New("gitsby.json contained no repos! :(")
	}
	return c, nil
}

func loadConfig() (GitsbyConfig, error) {
	envPath := os.Getenv("GITSBY_CONFIG")
	if envPath != "" {
		// Let's use the env var
		return readConfig(envPath)
	}
	// Check ~/gitsby/gitsby.json (default config)
	var usr, _ = user.Current()
	defaultPath := path.Join(usr.HomeDir, "gitsby", "gitsby.json")
	return readConfig(defaultPath)
}

func Setup() error {
	// Load our config file.
	if conf, loadErr := loadConfig(); loadErr != nil {
		return loadErr
	} else {
		Config = conf
	}
	// Set up our web.go server
	Server = web.NewServer()
	return nil
}
