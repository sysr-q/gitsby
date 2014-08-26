package main

import (
	"flag"
	"fmt"
	"github.com/hoisie/web"
	"io/ioutil"
	"log"
	"os/user"
	"path"
)

var config *Config
var server *web.Server

// Ignoring error here, I know. #dealwithit
var usr, _ = user.Current()

func gitsbyFolder(bits ...string) string {
	return path.Join(append([]string{usr.HomeDir, "gitsby"}, bits...)...)
}

var (
	host       = flag.String("host", "0.0.0.0", "host to bind web.go to")
	port       = flag.Int("port", 9999, "port to bind web.go to")
	configFile = flag.String("config", gitsbyFolder("gitsby.json"), "Gitsby config file")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime)
	flag.Parse()

	bindTo := fmt.Sprintf("%s:%d", *host, *port)

	if c, err := ParseConfig(*configFile); err == nil {
		config = c
	} else {
		log.Fatal(err)
	}

	// Yeah, no ternary operator leads to this.
	var plural string
	if len(config.Repos) != 1 {
		plural = "s"
	}
	log.Printf("The Great Gitsby is sending invites to %d repo%s.",
		len(config.Repos),
		plural)

	for _, repo := range config.Repos {
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

	server = web.NewServer()
	server.SetLogger(log.New(ioutil.Discard, "", 0))
	server.Post("/github", GitHub)

	log.Printf("The party is here: %s\n", bindTo)
	server.Run(bindTo)
}
