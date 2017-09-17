package main

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
)

type LinkFinder struct {
}

func (lf LinkFinder) FindDirectLink(filename string) (string, error) {
	abs_filename, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}
	r, err := exec.Command("ls", "-l", filepath.Dir(abs_filename)).Output()

	if err != nil {
		return "", err
	}

	lib_base_name := filepath.Base(abs_filename)

	for _, line := range strings.Split(string(r), "\n") {
		fields := strings.Fields(line)
		n := len(fields)

		if n > 3 && fields[0][0] == 'l' && fields[n-2] == "->" && fields[n-3] == lib_base_name {
			linkedFile := fields[n-1]
			if !filepath.IsAbs(linkedFile) {
				linkedFile = filepath.Join(filepath.Dir(abs_filename), linkedFile)
			}

			log.Infof("%s is symbol link of %s", filename, linkedFile)
			return linkedFile, nil
		}
	}
	return "", errors.New("no link found")
}

func (lf LinkFinder) IsSymbolLink(filename string) bool {
	_, err := lf.FindDirectLink(filename)
	return err == nil
}

//
// return a tuple (symbol link, linked file, full name, error )
//
func (lf LinkFinder) FindRealName(filename string) (string, string, string, error) {

	abs_filename, err := filepath.Abs(filename)

	if err != nil {
		return "", "", "", err
	}

	cur_filename := abs_filename
	for {
		dir := filepath.Dir(cur_filename)
		if dir == cur_filename {
			break
		}

		link, err := lf.FindDirectLink(dir)
		if err == nil && link != dir {
			return dir, link, strings.Replace(abs_filename, dir, link, 1), nil
		}
		cur_filename = dir
	}

	return "", "", "", errors.New("not a link")
}
