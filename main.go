package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hoisie/web"
	"io/ioutil"
	"log"
	"os/user"
	"path"
)

type configRepo struct {
	Url       string `json:"url"`
	Directory string `json:"directory"`
	Hidden    bool   `json:"hidden"`
	Type      string `json:"type"`
}

type Config struct {
	Repos       []Repo
	ReposActive map[string]Repo
}

var config Config
var server *web.Server

// Ignoring error here, I know. #dealwithit
var usr, _ = user.Current()

func gitsbyFolder(bits ...string) string {
	return path.Join(append([]string{usr.HomeDir, "gitsby"}, bits...)...)
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime)

	configDefault := gitsbyFolder("gitsby.json")

	host := flag.String("host", "0.0.0.0", "host to bind web.go to")
	port := flag.Int("port", 9999, "port to bind web.go to")
	configFile := flag.String("config", configDefault, "Gitsby config file")
	flag.Parse()
	bindTo := fmt.Sprintf("%s:%d", *host, *port)

	configData, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	var reps struct {
		Repos []configRepo `json:"repos"`
	}
	if err = json.Unmarshal(configData, &reps); err != nil {
		log.Fatal(err)
	}

	// Yeah, no ternary operator leads to this.
	var plural string
	if len(reps.Repos) != 1 {
		plural = "s"
	}
	log.Printf("The Great Gitsby is sending invites to %d repo%s.",
		len(reps.Repos),
		plural)

	var conf Config
	conf.Repos = []Repo{}
	conf.ReposActive = make(map[string]Repo)
	for _, crepo := range reps.Repos {
		// TODO: support SVN/CVS/$other?
		repo := Git{
			Url:       crepo.Url,
			directory: crepo.Directory,
			Hidden:    crepo.Hidden,
		}
		owner, name := repo.Name()
		conf.Repos = append(conf.Repos, repo)
		conf.ReposActive[owner+"/"+name] = repo
		var ok bool
		var err error
		if !repo.Exists() {
			repo.Log("doesn't exist, syncing!")
			ok, err = repo.Clone()
		} else {
			ok, err = repo.Pull()
		}

		if err == nil && ok {
			repo.Deploy()
		} else {
			repo.Log("unable to sync: %s\n", err)
		}
	}
	config = conf

	server = web.NewServer()
	server.SetLogger(log.New(ioutil.Discard, "", 0))
	server.Post("/github", GitHub)

	log.Printf("The party is here: %s\n", bindTo)
	server.Run(bindTo)
}
