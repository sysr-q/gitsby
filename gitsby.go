package main

import (
	"fmt"
	"log"
	// TODO: pass arounds flags instead of env vars.
	"flag"
	_ "html/template"
	_ "github.com/plausibility/gitsby/git"
	"github.com/plausibility/gitsby/util"
	"github.com/plausibility/gitsby/server"
	"github.com/plausibility/gitsby/github"
	"github.com/plausibility/gitsby/bitbucket"
)

func home() string {
	// TODO: render a template detailing what repos are being hooked,
	// (unless they're hidden), and wtf gitsby is.
	return ""
}

func main() {
	log.SetFlags(log.Ldate|log.Ltime)

	defaultConfig := util.GitsbyFolder("gitsby.json")

	host := flag.String("host", "0.0.0.0", "host to bind web.go to")
	port := flag.Int("port", 9999, "port to bind web.go to")
	config := flag.String("config", defaultConfig, "Gitsby config file")
	flag.Parse()

	bindTo := fmt.Sprintf("%s:%d", *host, *port)

	if setupErr := server.Setup(config); setupErr != nil {
		panic(setupErr)
	}

	log.Println("The Great Gitsby is throwing a party!")

	util.Infof("Preparing %d repo(s) for sync", len(server.Config.Repos))
	for _, repo := range server.Config.Repos {
		util.Infof("Checking %s", repo.Path())
		if !repo.Exists() {
			util.Info("Doesn't exist, syncing!")
			repo.Clone()
		}
	}

	if server.Config.Landing {
		server.Server.Get("/", home)
	}
	server.Server.Post("/github", github.Hook)
	server.Server.Post("/bitbucket", bitbucket.Hook)

	server.Server.Run(bindTo)
}
