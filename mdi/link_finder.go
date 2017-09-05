package main

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
)

type LinkFinder struct {
}

func (lf LinkFinder) FindLink(lib string, linkCallback func(link string)) {
	cur := lib
	for {
		link, err := lf.findDirectLink(cur)

		if err != nil {
			break
		}
		if link != cur {
			linkCallback(link)
			cur = link
		} else {
			break
		}
	}
}

func (lf LinkFinder) findDirectLink(lib string) (string, error) {
	r, err := exec.Command("ls", "-l", lib).Output()

	if err != nil {
		return "", err
	}

	line := string(r)

	fields := strings.Fields(line)
	n := len(fields)

	if n > 3 && fields[n-2] == "->" {

		if !filepath.IsAbs(fields[n-1]) {

			i := strings.LastIndex(lib, "/")

			if i > 0 {
				a := make([]string, 0)
				a = append(a, lib[0:i])
				a = append(a, fields[n-1])
				return strings.Join(a, "/"), nil
			}
		} else {
			return fields[n-1], nil
		}
	}
	return "", errors.New("no link found")
}
