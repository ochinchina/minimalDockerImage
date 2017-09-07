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
	includes []string
	excludes []string
}

func NewImageConfigWithFile(configFile string) (*ImageConfig, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	return NewImageConfig(f)
}

func NewImageConfig(reader io.Reader) (*ImageConfig, error) {
	config := &ImageConfig{includes: make([]string, 0), excludes: make([]string, 0)}

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
				config.excludes = append(config.excludes, line[1:])
			} else {
				config.includes = append(config.includes, line)
			}
		}
	}
	return config, nil
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
	for _, file := range ic.includes {
		file, err := filepath.Abs(file)
		if err != nil {
			continue
		}
		if ic.isWildcastFile(file) {
			result = append(result, path.Dir(file))
		} else {
			result = append(result, file)
		}
	}
	return result
}

// check if a file is in the include
//
func (ic *ImageConfig) inInclude(filename string) bool {
	return ic.inFileList(filename, ic.includes)

}

// check if a file is in the include
//
func (ic *ImageConfig) inExclude(filename string) bool {
	return ic.inFileList(filename, ic.excludes)
}

func (ic *ImageConfig) inFileList(filename string, files []string) bool {
	for _, file := range files {
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
