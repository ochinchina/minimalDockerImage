package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DependencyFinder struct {
	config     *ImageConfig
	result     []string
	linkFinder LinkFinder
}

func NewDependencyFinder(config *ImageConfig) *DependencyFinder {
	return &DependencyFinder{config: config, result: make([]string, 0)}
}

func (df *DependencyFinder) FindDependencies() []string {
	files := df.config.getAllIncludes()
	for len(files) > 0 {
		file := files[0]
		files = files[1:]
		//if it is in exclude list
		if df.config.inExclude(file) {
			continue
		}

		//if it is a link, find all the links
		df.linkFinder.FindLink(file, func(link string) {
			files = append(files, link)
		})

		//find all the dependencies
		if IsExecutable(file) {
			df.result = append(df.result, file)
			df.findDependencies(file, func(depLib string) {
				files = append(files, depLib)
			})
		} else if IsDir(file) {
			df.listFiles(file, func(f string) {
				files = append(files, f)
			})
		} else if df.config.inInclude(file) {
			df.result = append(df.result, file)
		}
	}

	return df.result
}
func (df *DependencyFinder) findDependencies(app string, depCallback func(depLib string)) {
	r, err := exec.Command("ldd", app).Output()

	if err != nil {
		return
	}

	lines := strings.Split(string(r), "\n")
	linkFinder := LinkFinder{}

	for _, line := range lines {
		fields := strings.Fields(line)
		dep := ""
		if len(fields) == 2 {
			if fields[1][0] == '(' {
				dep = fields[0]
			}
		} else if len(fields) == 4 {
			dep = fields[2]
		}
		if dep != "" {
			depCallback(dep)
		}

		linkFinder.FindLink(dep, depCallback)
	}
}
func (finder *DependencyFinder) listFiles(dir string, fileFoundCallback func(file string)) {
	filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileFoundCallback(file)
		return nil
	})
}
