package main

type Repo interface {
	Exists() bool
	Clone() (bool, error)
	Pull() (bool, error)
	Deploy() (bool, error)
	Name() (string, string) // owner, repo
	Path() string
}
