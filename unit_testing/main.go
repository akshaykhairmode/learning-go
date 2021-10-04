package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

}

func createFile(fileName string) (string, error) {

	if fileName == "" {
		return "", fmt.Errorf("Empty Path Given")
	}

	if !filepath.IsAbs(fileName) {
		abspath, err := filepath.Abs(fileName)
		if err != nil {
			return "", err
		}

		fileName = abspath
	}

	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		return "", err
	}

	return fileName, nil

}
