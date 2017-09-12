package main

import (
	"path/filepath"
	"strings"
)

type LinkManager struct {
	links    map[string]string
	notLinks map[string]string
	lf       LinkFinder
}

func NewLinkManager() *LinkManager {
	return &LinkManager{links: make(map[string]string),
		notLinks: make(map[string]string)}
}

func (lm *LinkManager) IsLink(filename string) bool {
	return lm.lf.IsSymbolLink(filename)
}

func (lm *LinkManager) FindRealName(filename string) (symbolLink string, realName string, err error) {
	abs_filename, err := filepath.Abs(filename)
	if err != nil {
		return
	}
	for symbol, f := range lm.links {
		if strings.HasPrefix(abs_filename, symbol) && abs_filename[len(symbol)] == filepath.Separator {
			symbolLink = symbol
			realName = strings.Replace(abs_filename, symbolLink, f, 1)
			return
		}
	}

	linkedName := ""
	symbolLink, linkedName, realName, err = lm.lf.FindRealName(abs_filename)
	if err == nil {
		lm.links[symbolLink] = linkedName
	} else {
		cur_file := abs_filename
		for {
			dir := filepath.Dir(cur_file)
			if dir == cur_file {
				break
			}

			lm.notLinks[dir] = dir
			cur_file = dir
		}
	}
	return
}

func (lm *LinkManager) mustNotSymbolLink(filename string) bool {
	cur_file := filename
	for {
		dir := filepath.Dir(cur_file)
		if _, ok := lm.notLinks[dir]; ok {
			return true
		}
		if dir == cur_file {
			break
		}
		cur_file = dir
	}
	return false
}
