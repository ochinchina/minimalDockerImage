package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// this class stores the found dependecies
//
//
type DependencyList struct {
	deps   map[string]string
	config *ImageConfig
}

func NewDependencyList(config *ImageConfig) *DependencyList {
	return &DependencyList{deps: make(map[string]string),
		config: config}
}

// Append a file to the list. If itself or its parent is already in the list, return false
// otherwise return true
//
func (dl DependencyList) Append(dep string) bool {
	if dl.config.inExclude(dep) {
		return false
	}
	for lib := range dl.deps {
		if IsParentDir(lib, dep) || !Exist(dep) {
			return false
		}
	}
	if _, ok := dl.deps[dep]; !ok {
		dl.deps[dep] = dep
		return true
	}
	return false
}

func (dl DependencyList) Contains(dep string) bool {
	_, ok := dl.deps[dep]
	return ok
}

// iterate the dependencies in the list
//
//
func (dl DependencyList) ForEach(depProcCallback func(dep string)) {
	tmpDeps := make([]string, 0)
	for k := range dl.deps {
		tmpDeps = append(tmpDeps, k)
	}
	sort.StringSlice(tmpDeps).Sort()
	for _, dep := range tmpDeps {
		depProcCallback(dep)
	}
}

func (dl *DependencyList) isLinkRetrieved(link string) bool {
	//check to see if any child in the link is retrieved
	for dep := range dl.deps {
		if IsParentDir(link, dep) {
			return true
		}
	}
	return false
}

type DependencyFinder struct {
	config     *ImageConfig
	result     *DependencyList
	linkFinder LinkFinder
	lm         *LinkManager
}

func NewDependencyFinder(config *ImageConfig) *DependencyFinder {
	return &DependencyFinder{config: config,
		result: NewDependencyList(config),
		lm:     NewLinkManager()}
}

func (df *DependencyFinder) FindDependencies() *DependencyList {
	pendingFiles := NewDependencyPendingFiles(df.config.getAllIncludes())

	for !pendingFiles.IsEmpty() {
		file, _ := pendingFiles.Take()
		log.Infof("Find dependency of %s", file)

		//if it contains symbol link
		symbolLink, realName, err := df.lm.FindRealName(file)
		if err == nil {
			pendingFiles.Add(realName)
			df.result.Append(symbolLink)
			continue
		}

		//if fail to add the file to the list
		if !df.result.Append(file) {
			continue
		}

		//if it is a link, find the direct link
		link, err := df.linkFinder.FindDirectLink(file)

		if err == nil {
			if !df.result.isLinkRetrieved(link) {
				pendingFiles.Add(link)
			}
			continue
		}

		//find all the dependencies
		if IsDir(file) {
			df.findDirDependencies(file, func(depLib string) {
				pendingFiles.Add(depLib)
			})
		} else if IsExecutable(file) {
			df.findDirectDependLibs(file, func(depLib string) {
				pendingFiles.Add(depLib)
			})
		}
	}
	return df.result
}

// Find the dependencies (out of the dir) of a directory
//
// Return: all the dependent libraries out of the directory
func (df *DependencyFinder) findDirDependencies(dir string, depCallback func(depLib string)) {
	lf := LinkFinder{}

	df.listFiles(dir, func(file string) {
		realFile, err := lf.FindDirectLink(file)
		if err == nil && !IsParentDir(dir, realFile) {
			depCallback(realFile)
		} else if IsExecutable(file) {
			df.findDirectDependLibs(file, func(depLib string) {
				if !IsParentDir(dir, depLib) { //if the dependent library is not under the dir
					depCallback(depLib)
				}
			})
		}
	})
}

// find the direct dependencies of a executable binary application or library
// Arguments:
//	app: the executable application or library
//  depCallback: the user provided callback to receive the depepdent libraries
func (df *DependencyFinder) findDirectDependLibs(app string, depCallback func(depLib string)) {
	r, err := exec.Command("ldd", app).Output()

	if err != nil {
		return
	}

	lines := strings.Split(string(r), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		dep := ""
		switch len(fields) {
		case 2:
			if fields[1][0] == '(' {
				dep = fields[0]
			}
		case 4:
			dep = fields[2]
		}

		if filepath.IsAbs(dep) {
			depCallback(dep)
		} else if dep != "" {
			dep = filepath.Join(filepath.Dir(app), dep)
			abs_dep, err := filepath.Abs(dep)
			if err == nil {
				depCallback(abs_dep)
			}
		}

	}
}
func (finder *DependencyFinder) listFiles(dir string, fileFoundCallback func(file string)) {
	abs_dir, err := filepath.Abs(dir)
	if err != nil {
		return
	}
	filepath.Walk(abs_dir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//notify a file is found
		abs_file, err := filepath.Abs(file)
		if err != nil {
			return err
		}
		if abs_dir != abs_file {
			fileFoundCallback(abs_file)
		}
		return nil
	})
}
