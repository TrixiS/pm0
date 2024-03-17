package utils

import (
	"os"
	"path"
)

const PM0Dirname = ".pm0"

func GetPM0Dirpath() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	pm0Dirpath := path.Join(homeDir, PM0Dirname)

	if err := os.MkdirAll(pm0Dirpath, 0777); err != nil {
		return "", err
	}

	return pm0Dirpath, nil
}
