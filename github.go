package main

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/web"
	"io/ioutil"
	"log"
)

type ghRepository struct {
	FullName string `json:"full_name"`
	Name     string `json:"name"`
	Url      string `json:"url"`
}

type ghPayload struct {
	After      string       `json:"after"`
	Repository ghRepository `json:"repository"`
}

func GitHub(ctx *web.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Printf("Unable to parse request body: %s\n", err)
		ctx.Abort(500, fmt.Sprintf("Unable to parse request body: %s", err.Error()))
		return
	}

	var p ghPayload
	if err := json.Unmarshal(body, &p); err != nil {
		log.Printf("Received unparseable payload: %s\n", err)
		ctx.Abort(500, fmt.Sprintf("Unable to parse payload: %s", err.Error()))
		return
	}

	// GitHub sends: owner/repo
	fullName := p.Repository.FullName

	repo, ok := config.ReposActive[fullName]
	if !ok {
		log.Printf("No repo found for: '%s', ignoring!\n", fullName)
		ctx.Abort(500, fmt.Sprintf("No repo found for: '%s'", fullName))
		return
	}

	log.Printf("Successful payload received for: '%s'!\n", fullName)
	// goroutine to cheat around `git pull` potentially blocking HTTP resp.
	go func() {
		if ok, err := repo.Pull(); err == nil && ok {
			// Only deploy if we pulled successfully.
			repo.Deploy()
		}
	}()
}
