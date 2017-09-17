package main

import (
	"errors"
)

type DependencyPendingFiles struct {
	pendingFiles []string
	allFiles     map[string]string
}

func NewDependencyPendingFiles(pendingFiles []string) *DependencyPendingFiles {
	dpf := &DependencyPendingFiles{pendingFiles: pendingFiles,
		allFiles: make(map[string]string)}
	if dpf.pendingFiles == nil {
		dpf.pendingFiles = make([]string, 0)
	}
	return dpf
}

func (dpf *DependencyPendingFiles) Add(filename string) bool {
	if dpf.contains(filename) ||
		dpf.isChildrenAdded(filename) ||
		dpf.isParentAdded(filename) {
		return false
	}

	dpf.pendingFiles = append(dpf.pendingFiles, filename)
	dpf.allFiles[filename] = filename
	return true
}

func (dpf *DependencyPendingFiles) isParentAdded(filename string) bool {
	for k := range dpf.allFiles {
		if IsParentDir(k, filename) {
			log.Warnf("the parent of %s is already added", filename)
			return true
		}
	}
	return false
}

func (dpf *DependencyPendingFiles) isChildrenAdded(filename string) bool {
	for k := range dpf.allFiles {
		if IsParentDir(filename, k) {
			log.Warnf("the child of %s is already added", filename)
			return true
		}
	}
	return false
}

func (dpf *DependencyPendingFiles) contains(filename string) bool {
	if _, ok := dpf.allFiles[filename]; ok {
		return true
	}
	return false
}

func (dpf *DependencyPendingFiles) Take() (string, error) {
	if len(dpf.pendingFiles) <= 0 {
		return "", errors.New("No pending files")
	}
	filename := dpf.pendingFiles[0]
	dpf.pendingFiles = dpf.pendingFiles[1:]
	return filename, nil
}

func (dpf *DependencyPendingFiles) IsEmpty() bool {
	return len(dpf.pendingFiles) <= 0
}
