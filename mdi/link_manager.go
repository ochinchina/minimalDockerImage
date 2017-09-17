package main

import (
	"errors"
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
	if lm.notSymbolLink(abs_filename) {
		return "", "", errors.New("Not a symbol")
	}
	symbolLink, realName, err = lm.findRealNameFromCache(abs_filename)
	if err != nil {
		return
	}

	linkedName := ""
	symbolLink, linkedName, realName, err = lm.lf.FindRealName(abs_filename)
	if err == nil {
		lm.links[symbolLink] = linkedName
	} else {
		lm.updateNotSymbol(abs_filename)
	}
	return
}

func (lm *LinkManager) updateNotSymbol(filename string) {
	cur_file := filename
	for {
		dir := filepath.Dir(cur_file)
		if dir == cur_file {
			break
		}

		lm.notLinks[dir] = dir
		cur_file = dir
	}
}
func (lm *LinkManager) findRealNameFromCache(filename string) (string, string, error) {
	for symbol, f := range lm.links {
		if strings.HasPrefix(filename, symbol) && filename[len(symbol)] == filepath.Separator {
			return symbol, strings.Replace(filename, symbol, f, 1), nil
		}
	}
	return "", "", errors.New("not a symbol")
}

func (lm *LinkManager) notSymbolLink(filename string) bool {
	dir := filepath.Dir(filename)
	if _, ok := lm.notLinks[dir]; ok {
		return true
	}
	return false
}
