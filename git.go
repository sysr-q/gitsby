package main

import (
	"os"
	"log"
	"regexp"
)

type Git struct {
	Url       string
	Directory string
	Hidden    bool
}

func (g Git) Exists() bool {
	return false
}

func (g Git) Clone() (bool, error) {
	return false, nil
}

func (g Git) Pull() (bool, error) {
	return false, nil
}

func (g Git) Deploy() (bool, error) {
	return false, nil
}

var metadataRegex = regexp.MustCompile(`[:/](?P<owner>\w+)/(?P<repo>.+?)(?:\.git)?$`)

func (g Git) Name() (string, string) {
	return RepoName(g.Url)
}

func RepoName(url string) (string, string) {
	if ms := metadataRegex.FindStringSubmatch(url); ms != nil {
		return ms[1], ms[2]
	}
	log.Printf("Shouldn't be here: git.RepoName(%q)\n", url)
	return "", ""
}

var homeRegex = regexp.MustCompile(`^~/`)
func (g Git) Path() string {
	if g.Directory != "" {
		if homeRegex.MatchString(g.Directory) {
			return homeRegex.ReplaceAllString(g.Directory, usr.HomeDir + string(os.PathSeparator))
		}
		return g.Directory
	}
	if owner, name := g.Name(); owner != "" && name != "" {
		// ~/gitsby/{{ owner }}/{{ name }}
		return gitsbyFolder(owner, name)
	}
	// TODO: handle this decently.
	return gitsbyFolder("damage-control")
}
