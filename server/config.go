package server

import "github.com/plausibility/gitsby/git"

type GitsbyConfig struct {
	Landing bool `json:"landing"`
	Repos []git.Repo `json:"repos"`
}
