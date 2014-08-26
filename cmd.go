package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
)

type CommandInfo struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

type Command struct {
	Path           string
	Info           CommandInfo
	Stdout, Stderr *bytes.Buffer
	Done           chan interface{}
	Error          error
	silent         bool
}

func (c *Command) Execute() {
	if !c.silent {
		fmt.Println("Executing:", c)
	}
	go func() {
		cmd := exec.Command(c.Info.Cmd, c.Info.Args...)
		cmd.Dir = c.Path
		cmd.Stdout = c.Stdout
		cmd.Stderr = c.Stderr
		c.Error = cmd.Run()
		c.Done <- nil
	}()
}

func NewCommand(cmdPath []string, info CommandInfo) *Command {
	return &Command{
		Path:   path.Join(cmdPath...),
		Info:   info,
		Stdout: new(bytes.Buffer),
		Stderr: new(bytes.Buffer),
		Done:   make(chan interface{}),
		silent: true,
	}
}
