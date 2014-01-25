package util

import "os"

// Two simple IO helpers.
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return false, err
}

func FolderExists(path string) (bool, error) {
	dir, err := os.Stat(path)
	if err == nil {
		return dir != nil && dir.IsDir(), nil
	}
	if os.IsNotExist(err) { return false, err }
	return false, err
}
