package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

type ImageConfig struct {
	includes []string
	excludes []string
}

func NewImageConfig(configFile string) (*ImageConfig, error) {
	config := &ImageConfig{includes: make([]string, 0), excludes: make([]string, 0)}
	f, err := os.Open(configFile)
	if err != nil {
		return config, err
	}

	r := bufio.NewReader(f)
	for {
		//read a line
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		//ignore the empty and comments line
		if len(line) <= 0 || line[0] == '#' {
			continue
		}

		if line[0] == '-' {
			config.excludes = append(config.excludes, line[1:])
		} else {
			config.includes = append(config.includes, line)
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
			pos := strings.Index(filename, file)
			if pos != -1 && file[pos+1] == os.PathSeparator {
				return true
			}
		} else if ic.isWildcastFile(file) && ic.matchPattern(filename, file) {
			return true
		}
	}
	return false
}

func (ic *ImageConfig) matchPattern(file string, file_pattern string) bool {
	regEx := ic.toRegEx(file_pattern)
	matched, err := regexp.MatchString(regEx, file)
	return matched && err == nil
}
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
