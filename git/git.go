package git

import (
	"errors"
	"regexp"
	"os/user"
	_ "os/exec"
	"github.com/plausibility/gitsby/util"
)

type Repo struct {
	Url string `json:"url"`
	Directory string `json:"directory"`
	Hidden bool `json:"hidden"`
}

func (r Repo) Exists() bool {
	exists, _ := util.DirectoryExists(r.Path())
	return exists
}

func (r Repo) Clone() {
	util.Successf("Clone: '%s' to '%s'", r.Url, r.Path())
}

func (r Repo) Pull() {
	util.Successf("Pull: '%s' in '%s'", r.Url, r.Path())
}

func (r Repo) Deploy() {
	util.Successf("Deploy: '%s' in '%s'", r.Url, r.Path())
}

var metadataRegex = regexp.MustCompile(`[:/](?:\w+)/(?P<repo>[^.]+)`)

func (r Repo) Name() string {
	if ms := metadataRegex.FindStringSubmatch(r.Url); ms != nil {
		return ms[1]
	}
	return nil
}

var usr, _ = user.Current()
var homeRegex = regexp.MustCompile(`^~/`)

func (r Repo) Path() string {
	if r.Directory != "" {
		if homeRegex.MatchString(r.Directory) {
			return homeRegex.ReplaceAllString(r.Directory, usr.HomeDir + "/")
		}
		return r.Directory
	}
	if name := r.Name(); name != nil {
		// ~/gitsby/{{ repo_name }}
		return util.GitsbyFolder(name)
	}
	panic(errors.New("OH GOD"))
}
