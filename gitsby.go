package main

import (
	"log"
	// TODO: pass arounds flags instead of env vars.
	_ "flag"
	_ "path"
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
	// TODO: use flag instead of hard coding.
	bindTo := "0.0.0.0:9999"
	if setupErr := server.Setup(); setupErr != nil {
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
