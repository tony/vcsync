package vcsync_test

import (
	"testing"

	"github.com/Masterminds/vcs"
	"github.com/tony/vcsync/vcsync"
)

func TestFindsVcsType(t *testing.T) {
	var vcsURLTests = []struct {
		url   string
		vtype vcs.Type
	}{
		{"git+https://github.com/tony/.dot-configs", vcs.Git},
		{"git+ssh://git@github.com/tony/roundup.git", vcs.Git},
		{"hg+http://foicica.com/hg/textadept", vcs.Hg},
		{"svn+http://svn.code.sf.net/p/docutils/code/trunk", vcs.Svn},
	}

	for _, tt := range vcsURLTests {
		u, err := vcsync.ParsePipURL(tt.url)

		if u.Vtype != tt.vtype {
			t.Errorf("vcs should resolve to %s, got: %v", tt.vtype, u.Vtype)
		}
		if err != nil {
			t.Error(err)
		}
	}

	var errTests = []struct {
		url string
		err error
	}{
		{"https://github.com/tony/.dot-configs", vcs.ErrCannotDetectVCS},
		{"lol+https://git@github.com/tony/roundup.git", vcs.ErrCannotDetectVCS},
	}
	for _, tt := range errTests {
		_, err := vcsync.ParsePipURL(tt.url)
		if err != tt.err {
			t.Errorf("url without vcs found should return %v, returned %v", tt.err, err)
		}
	}
}

func TestFindsRef(t *testing.T) {
	var vcsURLTests = []struct {
		url string
		ref string
	}{
		{"git+https://github.com/tony/.dot-configs@moo", "moo"},
		{"git+https://github.com/tony/.dot-configs", ""},
		{"git+ssh://git@github.com/tony/roundup.git@master", "master"},
		{"hg+http://foicica.com/hg/textadept@ha", "ha"},
		{"svn+http://svn.code.sf.net/p/docutils/code/trunk@2019", "2019"},
	}

	for _, tb := range vcsURLTests {
		u, err := vcsync.ParsePipURL(tb.url)
		if err != nil {
			t.Error(err)
		}

		if u.Ref != tb.ref {
			t.Errorf("vcs should resolve to %s, got: %v", tb.ref, u.Ref)
		}
	}
}

func TestFindsLocation(t *testing.T) {
	var vcsURLTests = []struct {
		url      string
		location string
	}{
		{"git+https://github.com/tony/.dot-configs@moo", "https://github.com/tony/.dot-configs"},
		{"git+ssh://git@github.com/tony/roundup.git@master", "ssh://github.com/tony/roundup.git"},
		{"hg+http://foicica.com/hg/textadept@ha", "http://foicica.com/hg/textadept"},
		{"svn+http://svn.code.sf.net/p/docutils/code/trunk@2019", "http://svn.code.sf.net/p/docutils/code/trunk"},
	}

	for _, tb := range vcsURLTests {
		u, err := vcsync.ParsePipURL(tb.url)
		if err != nil {
			t.Error(err)
		}

		if u.Location() != tb.location {
			t.Errorf("vcs should resolve to %s, got: %v", tb.location, u.Location)
		}
	}
}
