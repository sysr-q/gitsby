package main

import (
	"log"
	"encoding/json"
	"github.com/hoisie/web"
)

type ghRepository struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type ghPayload struct {
	After string `json:"after"`
	Repository ghRepository `json:"repository"`
}

func GitHub(ctx *web.Context) {
	var p ghPayload
	if err := json.Unmarshal([]byte(ctx.Params["payload"]), &p); err != nil {
		ctx.Abort(500, "")
		return
	}
	
	owner, name := RepoName(p.Repository.Url)
	if owner == "" || name == "" {
		log.Printf("Received unparseable payload for: '%s/%s', ignoring!\n", owner, name)
		ctx.Abort(500, "")
		return
	}

	repo, ok := config.ReposActive[owner+"/"+name]
	if !ok {
		log.Printf("No repo found for: '%s/%s', ignoring!\n", owner, name)
		ctx.Abort(500, "")
		return
	}

	log.Printf("Successful payload received for: '%s/%s'!\n", owner, name)
	// goroutine to cheat around `git pull` potentially blocking HTTP resp.
	go func() {
		if ok, err := repo.Pull(); err == nil && ok {
			// Only deploy if we pulled successfully.
			repo.Deploy()
		}
	}();
}
