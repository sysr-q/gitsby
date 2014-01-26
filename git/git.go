package git

import (
	"path"
	"bytes"
	"errors"
	"regexp"
	"strings"
	"os/user"
	"os/exec"
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

func (r Repo) Clone() (bool, error) {
	owner, name := r.Name()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	// git clone git@github.com:owner/repo.git owner/repo
	cmd := exec.Command("git", "clone", r.Url, path.Join(owner, name))
	cmd.Dir = path.Join(r.Path(), "..", "..")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		util.Failuref("[%s/%s] failed to clone! Here's why: (stderr)", owner, name)
		util.Red(strings.TrimSuffix(stderr.String(), "\n"))
		return false, err
	}
	util.Successf("[%s/%s] successfully cloned to: %s", owner, name, r.Path())
	util.Blue(strings.TrimSuffix(stdout.String(), "\n"))
	return true, nil
}

func (r Repo) Pull() (bool, error) {
	owner, name := r.Name()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	// git pull origin
	cmd := exec.Command("git", "pull", "origin")
	cmd.Dir = r.Path()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		util.Failuref("[%s/%s] failed to pull! Here's why:", owner, name)
		util.Red(strings.TrimSuffix(stderr.String(), "\n"))
		return false, err
	}
	util.Successf("[%s/%s] synced", owner, name)
	util.Blue(strings.TrimSuffix(stdout.String(), "\n"))
	return true, nil
}

func (r Repo) Deploy() (bool, error) {
	owner, name := r.Name()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	// make autodeploy
	cmd := exec.Command("make", "autodeploy")
	cmd.Dir = r.Path()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		util.Failuref("[%s/%s] failed to 'make autodeploy'! Here's why:", owner, name)
		util.Red(strings.TrimSuffix(stderr.String(), "\n"))
		return false, err
	}
	util.Successf("[%s/%s] autodeployed!", owner, name)
	util.Blue(strings.TrimSuffix(stdout.String(), "\n"))
	return true, nil
}

var metadataRegex = regexp.MustCompile(`[:/](?P<owner>\w+)/(?P<repo>.+?)(?:\.git)?$`)

func (r Repo) Name() (string, string) {
	return RepoName(r.Url)
}

func RepoName(url string) (string, string) {
	if ms := metadataRegex.FindStringSubmatch(url); ms != nil {
		return ms[1], ms[2]
	}
	return "", ""
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
	if owner, name := r.Name(); owner != "" && name != "" {
		// ~/gitsby/{{ owner }}/{{ name }}
		return util.GitsbyFolder(owner, name)
	}
	panic(errors.New("Unable to generate Path() for repo: " + r.Url))
}
