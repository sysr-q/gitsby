package main

import (
	"flag"
	"fmt"
	"net/http"
	"log"
	"os/user"
	"path"
)

var config *Config

// Ignoring error here, I know. #dealwithit
var usr, _ = user.Current()

func gitsbyFolder(bits ...string) string {
	return path.Join(append([]string{usr.HomeDir, "gitsby"}, bits...)...)
}

var (
	host       = flag.String("host", "0.0.0.0", "host to bind net/http to")
	port       = flag.Int("port", 9999, "port to bind net/http to")
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

	log.Printf("The Great Gitsby is sending invites to %d repo(s).",
		len(config.Repos))

	for _, repo := range config.Repos {
		var (
			ok bool
			err error
		)

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

	log.Printf("The party is here: %s\n", bindTo)

	http.HandleFunc("/github", GitHub)
	log.Fatal(http.ListenAndServe(bindTo, nil))
}
