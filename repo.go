package main

type Repo interface {
	Exists() bool
	Clone() (bool, error)
	Pull() (bool, error)

	Deploy() (bool, error)

	// (owner, repo)
	Name() (string, string)
	Path() string

	Log(string, ...interface{})
}
