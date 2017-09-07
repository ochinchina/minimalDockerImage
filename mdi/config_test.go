package main

import (
	"strings"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	reader := strings.NewReader("/usr/bin/ls\n/usr/include/\n-/usr/include/boost\n-/usr/include/python*")
	config, err := NewImageConfig(reader)

	if err != nil {
		t.Fail()
	}

	if !config.inExclude("/usr/include/boost") || !config.inExclude("/usr/include/python2.7") {
		t.Fail()
	}
}
