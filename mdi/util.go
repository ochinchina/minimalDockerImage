package main

import "os"

func IsDir(file string) bool {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func IsExecutable(file string) bool {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return false
	}

	return fileInfo.Mode()&0111 != 0

}
