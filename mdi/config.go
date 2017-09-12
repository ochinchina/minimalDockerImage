package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type ImageConfig struct {
	includes map[string]bool
	excludes map[string]bool
}

func NewImageConfigWithFile(configFile string) (*ImageConfig, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	return NewImageConfig(f)
}

func NewImageConfig(reader io.Reader) (*ImageConfig, error) {
	config := &ImageConfig{includes: make(map[string]bool), excludes: make(map[string]bool)}

	//read whole content
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	//process line by line
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 0 && line[0] != '#' {
			if line[0] == '-' {
				config.addFile(line[1:], &config.excludes)
			} else {
				config.addFile(line, &config.includes)
			}
		}
	}
	return config, nil
}

func (ic *ImageConfig) AddInclude(filename string) error {
	return ic.addFile(filename, &ic.includes)

}

func (ic *ImageConfig) AddExclude(filename string) error {
	return ic.addFile(filename, &ic.excludes)
}

func (ic *ImageConfig) addFile(filename string, files *map[string]bool) error {
	abs_filename, err := filepath.Abs(filename)
	if err == nil {
		(*files)[abs_filename] = true
	}
	return nil
}

// Check if the file is a wildcast file or not
//
// Returns: true if the filename contains wildcast
func (ic *ImageConfig) isWildcastFile(filename string) bool {
	basename := path.Base(filename)

	return strings.IndexAny(basename, "*?") != -1
}

func (ic *ImageConfig) getAllIncludes() []string {
	result := make([]string, 0)
	links := make([]string, 0)
	lf := LinkFinder{}
	//the symbol link should be in the tail of the result list
	for file := range ic.includes {
		if ic.isWildcastFile(file) {
			file = path.Dir(file)
		}
		if lf.IsSymbolLink(file) {
			links = append(links, file)
		} else {
			result = append(result, file)
		}
	}
	for _, link := range links {
		result = append(result, link)
	}
	return result
}

// check if a file is in the include
//
func (ic *ImageConfig) inInclude(filename string) bool {
	return ic.inFileList(filename, &ic.includes)

}

// check if a file is in the include
//
func (ic *ImageConfig) inExclude(filename string) bool {
	return ic.inFileList(filename, &ic.excludes)
}

func (ic *ImageConfig) inFileList(filename string, files *map[string]bool) bool {
	for file := range *files {
		if file == filename {
			return true
		}
		if IsDir(file) {
			return IsParentDir(file, filename)
		} else if ic.isWildcastFile(file) && ic.matchPattern(filename, file) {
			return true
		}
	}
	return false
}

// check if the file matches the file_pattern or not
// Returns:
//  true if the file matches the file_pattern
func (ic *ImageConfig) matchPattern(file string, file_pattern string) bool {
	regEx := ic.toRegEx(file_pattern)
	matched, err := regexp.MatchString(regEx, file)
	return matched && err == nil
}

// convert the supported '*' and '?' to golang regex
// Arguments:
//	s the string maybe include '*' or '?'
// Returns:
// the converted golang regex
func (imageConfig *ImageConfig) toRegEx(s string) string {
	result := ""
	for _, ch := range s {
		if ch == '*' {
			result = result + ".*"
		} else if ch == '?' {
			result = result + "."
		} else if ch == '.' {
			result = result + "\\."
		} else {
			result = fmt.Sprintf("%s%c", result, ch)
		}
	}
	return result
}
