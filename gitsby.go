package main

import (
	"fmt"
	// TODO: pass arounds flags instead of env vars.
	_ "flag"
	_ "path"
	_ "html/template"
	"github.com/hoisie/web"
	"github.com/plausibility/gitsby/server"
)

func home() string {
	// TODO: render a template detailing what repos are being hooked,
	// (unless they're hidden), and wtf gitsby is. 
	return ""
}

func hook(ctx *web.Context) string {
	return ""
}

func main() {
	if setupErr := server.Setup(); setupErr != nil {
		panic(setupErr)
	}
	fmt.Println(server.Config)
	if server.Config.Landing {
		server.Server.Get("/", home)
	}
	server.Server.Post("/hook", hook)
	server.Server.Run("0.0.0.0:9999")
}
