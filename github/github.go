package github

import (
	"github.com/hoisie/web"
	"github.com/plausibility/gitsby/git"
	"github.com/plausibility/gitsby/util"
	"github.com/plausibility/gitsby/server"
)

func Hook(ctx *web.Context) {
	p := NewPayload([]byte(ctx.Params["payload"]))
	owner, name := git.RepoName(p.Url())
	if owner == "" || name == "" {
		util.Infof("Received unparseable payload: '%s/%s', ignoring!", owner, name)
		return
	}
	repo, ok := server.Config.ReposActive[owner + "/" + name]
	if !ok {
		util.Infof("Received payload for unknown repo: '%s/%s', ignoring!", owner, name)
		return
	}
	util.Info("Received payload for: '%s/%s'!", owner, name)
	// goroutines to cheat around `git pull` blocking HTTP.
	go func() {
		if ok, err := repo.Pull(); err == nil && ok {
			// Only deploy if we pull'd successfully.
			repo.Deploy()
		}
	}();
}
