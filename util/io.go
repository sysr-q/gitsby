package util

import (
	"os"
	"path"
	"os/user"
)

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return false, err
}

func DirectoryExists(path string) (bool, error) {
	dir, err := os.Stat(path)
	if err == nil {
		return dir != nil && dir.IsDir(), nil
	}
	if os.IsNotExist(err) { return false, err }
	return false, err
}

func GitsbyFolder(subs ...string) (string) {
	usr, _ := user.Current()
	base := path.Join(usr.HomeDir, "gitsby")
	for _, e := range subs {
		base = path.Join(base, e)
	}
	return base
}
