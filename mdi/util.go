package main

import (
	"os"
	"strings"
)

type FileInfo struct {
	fileInfo os.FileInfo
	err      error
}

func NewFileInfo(file string) *FileInfo {
	fileInfo, err := os.Stat(file)

	return &FileInfo{fileInfo: fileInfo, err: err}
}

func (fi *FileInfo) IsDir() bool {
	return fi.err == nil && fi.fileInfo.IsDir()
}

func (fi *FileInfo) Exist() bool {
	return fi.err == nil
}

func (fi *FileInfo) IsExecutable() bool {
	return fi.err == nil && fi.fileInfo.Mode().Perm()&0x111 != 0
}

func IsDir(file string) bool {
	return NewFileInfo(file).IsDir()
}

func IsExecutable(file string) bool {
	return NewFileInfo(file).IsExecutable()
}

func Exist(file string) bool {
	return NewFileInfo(file).Exist()
}

// check if the parent is the parent of dir
//
// return: true if the parent is the parent of dir

func IsParentDir(parent, dir string) bool {
	n := len(parent) - 1
	if parent[n] == os.PathSeparator {
		parent = parent[0:n]
		n--
	}
	if strings.HasPrefix(dir, parent) {
		t := dir[n+1:]
		return len(t) == 0 || t[0] == os.PathSeparator

	}
	return false
}
