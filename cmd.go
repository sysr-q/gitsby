package main

import (
	"fmt"
	"bytes"
	"os/exec"
	"path"
)

type Command struct {
	Path, Cmd      string
	Args           []string
	Stdout, Stderr *bytes.Buffer
	Done           chan interface{}
	Error          error
	silent bool
}

func (c *Command) Execute() {
	if !c.silent {
		fmt.Println("Executing:", c)
	}
	go func() {
		cmd := exec.Command(c.Cmd, c.Args...)
		cmd.Dir = c.Path
		cmd.Stdout = c.Stdout
		cmd.Stderr = c.Stderr
		c.Error = cmd.Run()
		c.Done <- nil
	}()
}

func NewCommand(cmdPath []string, cmd string, args []string) *Command {
	return &Command{
		Path:   path.Join(cmdPath...),
		Cmd:    cmd,
		Args:   args,
		Stdout: new(bytes.Buffer),
		Stderr: new(bytes.Buffer),
		Done:   make(chan interface{}),
		silent: true,
	}
}
