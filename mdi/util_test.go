package main

import (
	"testing"
)

func TestIsParentDir(t *testing.T) {
	dir := "/this/is/a/dir"
	parent := "/this/is"

	if !IsParentDir(parent, dir) {
		t.Fatal()
	}

	parent = "/this/is/a/"

	if !IsParentDir(parent, dir) {
		t.Fatal()
	}

	if !IsParentDir("/usr/lib64", "/usr/lib64/libexempi.so.3") {
		t.Fatal()
	}
}
