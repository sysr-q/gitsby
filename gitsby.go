package main

import (
	"fmt"
	_ "path"
	_ "html/template"
	"github.com/hoisie/web"
	"github.com/plausibility/gitsby/server"
)

func home() string {
	return "hello"
}

func hook(ctx *web.Context) string {
	// todo
	// tpl_path := path.Join()
	// tpl, err := template.ParseFiles(tpl_path)
	return ""
}

func main() {
	if setupErr := server.Setup(); setupErr != nil {
		panic(setupErr)
	}
	fmt.Println(server.Config)
	server.Server.Get("/", home)
	server.Server.Post("/hook", hook)
	server.Server.Run("0.0.0.0:9999")
}
