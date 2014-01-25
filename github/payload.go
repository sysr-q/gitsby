package github

import "encoding/json"

// NOTE: I left out some non-essential things from the structs.
// To be fair, we only _really_ need the URL; anything else is a plus!

// Doubles as an Owner/Committer
type Author struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Username string `json:"username"`
}

type Commit struct {
	Id string `json:"id"`
	Message string `json:"message"`
	Timestamp string `json:"timestamp"`
	Url string `json:"url"`
	Added []string `json:"added"`
	Removed []string `json:"removed"`
	Modified []string `json:"modified"`
	Author Author `json:"author"`
	Committer Author `json:"committer"`
}

type Repository struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Description string `json:"description"`
	Homepage string `json:"homepage"`
	Watchers int `json:"watchers"`
	Forks int `json:"forks"`
	Private bool `json:"private"`
	Owner Author `json:"owner"`
}

type Payload struct {
	Before string `json:"before"`
	After string `json:"after"`
	Ref string `json:"ref"`
	Commits []Commit `json:"commits"`
	Repository Repository `json:"repository"`
}

func NewPayload(data []byte) Payload {
	var payload Payload
	if err := json.Unmarshal(data, &payload); err != nil {
		panic(err)
	}
	return payload
}

// The only thing we really care about.
func (p *Payload) Url() string {
	return p.Repository.Url
}
