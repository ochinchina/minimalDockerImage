package main

import "testing"

func TestLinkFinder_FindDirectLink(t *testing.T) {
	lf := LinkFinder{}
	link, err := lf.FindDirectLink("/usr/bin/ld")
	if err != nil || link != "/etc/alternatives/ld" {
		t.Fail()
	}
}
