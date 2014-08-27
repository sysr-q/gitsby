package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func GitHub(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to parse request body: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to parse request body: %s", err.Error())
		return
	}

	var p ghPayload
	if err := json.Unmarshal(body, &p); err != nil {
		log.Printf("Received unparseable payload: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to parse payload: %s", err.Error())
		return
	}

	// GitHub sends: owner/repo
	fullName := p.Repository.FullName

	repo, ok := config.Repos[fullName]
	if !ok {
		log.Printf("No repo found for: %s, ignoring!\n", fullName)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "No repo found for: %s", fullName)
		return
	}

	log.Printf("Successful payload received for: %s!\n", fullName)
	w.Write([]byte("Success"))

	// goroutine to cheat around `git pull` potentially blocking HTTP resp.
	go func() {
		if ok, err := repo.Pull(); err == nil && ok {
			// Only deploy if we pulled successfully.
			repo.Deploy()
		}
	}()
}
