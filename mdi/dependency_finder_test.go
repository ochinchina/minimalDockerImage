package main

import (
	"strings"
	"testing"
)

func TestDependencyList_Contains(t *testing.T) {
	reader := strings.NewReader("/usr/lib64/libssl3.so")
	config, err := NewImageConfig(reader)
	if err != nil {
		t.Fatal("Fail to load the configure")
	}
	depList := NewDependencyList(config)
	if !depList.Append("/usr/lib64/libssl3.so") {
		t.Fatal("Fail to add lib")
	}
	if depList.Append("/usr/lib64/libssl3.so") {
		t.Fatal("lib is already added")
	}

	if !depList.Contains("/usr/lib64/libssl3.so") {
		t.Fatal("lib should be in the list")
	}
}
