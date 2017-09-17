package main

import (
	"testing"
)

func TestDependencyPendingFiles_Add(t *testing.T) {
	pendingFiles := NewDependencyPendingFiles(nil)
	if !pendingFiles.Add("/home/test/test1") {
		t.Fatal("fail to add file")
	}
	if pendingFiles.Add("/home/test/test1") {
		t.Fatal("same file can't add twice")
	}
}

func TestDependencyPendingFiles_Take(t *testing.T) {
	pendingFiles := NewDependencyPendingFiles(nil)
	pendingFiles.Add("/home/test/test1")
	pendingFiles.Add("/home/test/test2")
	if pendingFiles.IsEmpty() {
		t.Fatal("not empty")
	}
	f, _ := pendingFiles.Take()
	if f != "/home/test/test1" {
		t.Fatal("first added should be got at first")
	}

	f, _ = pendingFiles.Take()
	if f != "/home/test/test2" {
		t.Fatal("seocnd added shoud be got secondly")
	}

	if !pendingFiles.IsEmpty() {
		t.Fatal("all elements should be taken")
	}
}

func TestDependencyPendingFiles_FailToAddIfParentIsAdded(t *testing.T) {
	pendingFiles := NewDependencyPendingFiles(nil)
	if !pendingFiles.Add("/home/test") {
		t.Fail()
	}

	if pendingFiles.Add("/home/test/test2") {
		t.Fatalf("child should not be added because the parent is already added")
	}
}

func TestDependencyPendingFiles_FailToAddIfChildIsAdded(t *testing.T) {
	pendingFiles := NewDependencyPendingFiles(nil)
	pendingFiles.Add("/home/test/test1")
	if pendingFiles.Add("/home/test") {
		t.Fatalf("parent should not be added because the child is already added")
	}
}
