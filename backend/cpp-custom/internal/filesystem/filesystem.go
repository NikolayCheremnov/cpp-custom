package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadFileToString(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func CreateWithClosing(path string) (string, error) {
	path, f, err := CreateWithRelativePath(path)
	if err != nil {
		return "", err
	}
	f.Close()
	return path, nil
}

func CreateWithRelativePath(path string) (string, *os.File, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return Create(path)
}

func Create(path string) (string, *os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return "", nil, err
	}
	file, err := os.Create(path)
	return path, file, err
}
