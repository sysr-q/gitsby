package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
)

type Git struct {
	Url       string
	directory string
	Hidden    bool
	Silent    bool
}

func (g Git) Log(fmt string, in ...interface{}) {
	if g.Silent {
		return
	}
	owner, name := g.Name()
	in = append([]interface{}{owner, name}, in...)
	log.Printf("[%s/%s] "+fmt, in...)
}

func (g Git) Print(in ...interface{}) {
	if g.Silent {
		return
	}
	fmt.Print(in...)
}

func (g Git) Exists() bool {
	_, err := os.Stat(g.Path())
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (g Git) Clone() (bool, error) {
	var cwd []string
	var args []string
	if g.Directory() == "" {
		owner, name := g.Name()
		cwd = []string{g.Path(), "..", ".."}
		args = []string{"clone", g.Url, path.Join(owner, name)}
	} else {
		cwd = []string{g.Path(), ".."}
		args = []string{"clone", g.Url,  g.Directory()}
	}

	cmd := NewCommand(cwd, "git", args)
	cmd.Execute()
	<-cmd.Done

	if cmd.Error != nil {
		g.Log("failed to clone! Here's why: (stderr)")
		g.Print(cmd.Stderr.String())
		return false, cmd.Error
	}
	g.Log("successfully cloned to: %s (stdout)", g.Path())
	g.Print(cmd.Stdout.String())
	return true, nil
}

func (g Git) Pull() (bool, error) {
	cmd := NewCommand([]string{g.Path()}, "git", []string{"pull", "origin"})
	cmd.Execute()
	<-cmd.Done

	if cmd.Error != nil {
		g.Log("failed to pull! Here's why: (stderr)")
		g.Print(cmd.Stderr.String())
		return false, cmd.Error
	}
	g.Log("successfully synced: (stdout)")
	g.Print(cmd.Stdout.String())
	return true, nil
}

func (g Git) Deploy() (bool, error) {
	cmd := NewCommand([]string{g.Path()}, "make", []string{"autodeploy"})
	cmd.Execute()
	<-cmd.Done

	if cmd.Error != nil {
		g.Log("failed to deploy! Here's why: (stderr)")
		g.Print(cmd.Stderr.String())
		return false, cmd.Error
	}
	g.Log("successfully deployed: (stdout)")
	g.Print(cmd.Stdout.String())
	return true, nil
}

var metadataRegex = regexp.MustCompile(`[:/](?P<owner>\w+)(?:/(?P<repo>.+?))?(?:\.git)?$`)

func (g Git) Name() (string, string) {
	return RepoName(g.Url)
}

func RepoName(url string) (string, string) {
	if ms := metadataRegex.FindStringSubmatch(url); ms != nil {
		return ms[1], ms[2]
	}
	// TODO: handle this decently.
	log.Printf("Shouldn't be here: git.RepoName(%q) - damage control!\n", url)
	return "", ""
}

var homeRegex = regexp.MustCompile(`^~/`)

func (g Git) Path() string {
	dir := g.Directory()
	if dir != "" {
		return dir
	}
	if owner, name := g.Name(); owner != "" && name != "" {
		// ~/gitsby/{{ owner }}/{{ name }}
		return gitsbyFolder(owner, name)
	}
	// TODO: handle this decently.
	return gitsbyFolder("damage-control")
}

func (g Git) Directory() string {
	if homeRegex.MatchString(g.directory) {
		return homeRegex.ReplaceAllString(g.directory, usr.HomeDir+string(os.PathSeparator))
	}
	return g.directory
}
