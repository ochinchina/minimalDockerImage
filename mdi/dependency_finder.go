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
	deps map[string]bool
}

func NewDependencyList() *DependencyList {
	return &DependencyList{deps: make(map[string]bool)}
}

// Append a file to the list. If itself or its parent is already in the list, return false
// otherwise return true
//
func (dl DependencyList) Append(dep string) bool {
	for lib := range dl.deps {
		if IsParentDir(lib, dep) || !Exist(dep) {
			return false
		}
	}
	if _, ok := dl.deps[dep]; !ok {
		dl.deps[dep] = true
		return true
	}
	return false
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

type DependencyFinder struct {
	config     *ImageConfig
	result     *DependencyList
	linkFinder LinkFinder
}

func NewDependencyFinder(config *ImageConfig) *DependencyFinder {
	return &DependencyFinder{config: config,
		result: NewDependencyList()}
}

func (df *DependencyFinder) FindDependencies() *DependencyList {
	files := df.config.getAllIncludes()
	for len(files) > 0 {
		file := files[0]
		files = files[1:]
		//if it is in exclude list
		if df.config.inExclude(file) || !df.result.Append(file) {
			continue
		}

		//if it is a link, find all the links
		is_link := false
		df.linkFinder.FindLink(file, func(link string) {
			files = append(files, link)
			is_link = true
		})

		if is_link {
			continue
		}

		//find all the dependencies
		if IsExecutable(file) {
			df.findDirectDependLibs(file, func(depLib string) {
				files = append(files, depLib)
			})
		} else if IsDir(file) {
			df.findDirDependencies(file, func(depLib string) {
				files = append(files, depLib)
			})
		}
	}
	return df.result
}

// Find the dependencies (out of the dir) of a directory
//
// Return: all the dependent libraries out of the directory
func (df *DependencyFinder) findDirDependencies(dir string, depCallback func(depLib string)) {
	dirs := []string{dir}

	for len(dirs) > 0 {
		cur_dir := dirs[0]
		dirs = dirs[1:]

		df.listFiles(cur_dir, func(file string) {
			if IsDir(file) { //recursively find the dependency libraries
				dirs = append(dirs, file)
			} else if IsExecutable(file) {
				df.findDirectDependLibs(file, func(depLib string) {
					if !IsParentDir(dir, depLib) { //if the dependent library is not under the dir
						depCallback(depLib)
					}
				})
			}
		})
	}

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
		if dep != "" {
			depCallback(dep)
		}
	}
}
func (finder *DependencyFinder) listFiles(dir string, fileFoundCallback func(file string)) {
	filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//notify a file is found
		fileFoundCallback(file)
		return nil
	})
}
