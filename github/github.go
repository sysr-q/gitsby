package github

import (
	"github.com/hoisie/web"
	"github.com/plausibility/gitsby/git"
	"github.com/plausibility/gitsby/util"
	"github.com/plausibility/gitsby/server"
)

func Hook(ctx *web.Context) {
	p := NewPayload([]byte(ctx.Params["payload"]))
	name := git.RepoName(p.Url())
	if name == "" { return }
	repo, ok := server.Config.ReposActive[name]
	if !ok {
		util.Infof("Received payload for unknown repo: '%s', ignoring!", name)
		return
	}
	if ok, err := repo.Pull(); err == nil && ok {
		repo.Deploy()
	}
}
