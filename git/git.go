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
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("git", "clone", r.Url)
	cmd.Dir = path.Join(r.Path(), "..")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		util.Failuref("[%s] failed to clone! Here's why: (stderr)", r.Name())
		util.Red(strings.TrimSuffix(stderr.String(), "\n"))
		return false, err
	}
	util.Successf("[%s] successfully cloned to: %s", r.Name(), r.Path())
	util.Blue(strings.TrimSuffix(stdout.String(), "\n"))
	return true, nil
}

func (r Repo) Pull() (bool, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("git", "pull", "origin")
	cmd.Dir = r.Path()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		util.Failuref("[%s] failed to pull! Here's why:", r.Name())
		util.Red(strings.TrimSuffix(stderr.String(), "\n"))
		return false, err
	}
	util.Successf("[%s] synced", r.Name())
	util.Blue(strings.TrimSuffix(stdout.String(), "\n"))
	return true, nil
}

func (r Repo) Deploy() (bool, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("make", "autodeploy")
	cmd.Dir = r.Path()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		util.Failuref("[%s] failed to `make autodeploy`! Here's why:", r.Name())
		util.Red(strings.TrimSuffix(stderr.String(), "\n"))
		return false, err
	}
	util.Successf("[%s] autodeployed!", r.Name())
	util.Blue(strings.TrimSuffix(stdout.String(), "\n"))
	return true, nil
}

var metadataRegex = regexp.MustCompile(`[:/](?:\w+)/(?P<repo>[^.]+)`)

func (r Repo) Name() string {
	return Name(r.Url)
}

func Name(url string) string {
	if ms := metadataRegex.FindStringSubmatch(url); ms != nil {
		return ms[1]
	}
	return ""
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
	if name := r.Name(); name != "" {
		// ~/gitsby/{{ repo_name }}
		return util.GitsbyFolder(name)
	}
	panic(errors.New("Unable to generate Path() for repo: " + r.Url))
}
