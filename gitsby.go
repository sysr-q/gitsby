package main

import (
	"fmt"
	"log"
	// TODO: pass arounds flags instead of env vars.
	"flag"
	"os/user"
	"path"
	_ "html/template"
	_ "github.com/plausibility/gitsby/git"
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

	usr, _ := user.Current()
	defaultConfig := path.Join(usr.HomeDir, "gitsby", "gitsby.json")

	host := flag.String("host", "0.0.0.0", "host to bind web.go to")
	port := flag.Int("port", 9999, "port to bind web.go to")
	config := flag.String("config", defaultConfig, "Gitsby config file")
	flag.Parse()

	bindTo := fmt.Sprintf("%s:%d", *host, *port)

	if setupErr := server.Setup(config); setupErr != nil {
		panic(setupErr)
	}

	log.Printf("The Great Gitsby is throwing a party at: %s", bindTo)

	if server.Config.Landing {
		server.Server.Get("/", home)
	}
	server.Server.Post("/github", github.Hook)
	server.Server.Post("/bitbucket", bitbucket.Hook)

	server.Server.Run(bindTo)
}
