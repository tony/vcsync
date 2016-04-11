package vcsync_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Masterminds/vcs"
	"github.com/tony/vcsync/vcsync"
)

func TestFindsVcsType(t *testing.T) {
	var configTests = []struct {
		url   string
		vtype vcs.Type
	}{
		{"git+https://github.com/tony/.dot-configs", vcs.Git},
		{"git+ssh://git@github.com/tony/roundup.git", vcs.Git},
		{"hg+http://foicica.com/hg/textadept", vcs.Hg},
		{"svn+http://svn.code.sf.net/p/docutils/code/trunk", vcs.Svn},
	}

	for _, tt := range configTests {
		vcsinfo, err := vcsync.ParsePIPUrl(tt.url)

		if vcsinfo.VCSType != tt.vtype {
			t.Errorf("vcs should resolve to %s, got: %v", tt.vtype, vcsinfo.VCSType)
		}
		if err != nil {
			t.Error(err)
		}
	}
}

func TestFindsBranch(t *testing.T) {
	var configTests = []struct {
		url    string
		branch string
	}{
		{"git+https://github.com/tony/.dot-configs@moo", "moo"},
		{"git+ssh://git@github.com/tony/roundup.git@master", "master"},
		{"hg+http://foicica.com/hg/textadept@ha", "ha"},
		{"svn+http://svn.code.sf.net/p/docutils/code/trunk@2019", "2019"},
	}

	for _, tb := range configTests {
		vcsinfo, err := vcsync.ParsePIPUrl(tb.url)

		if vcsinfo.Branch != tb.branch {
			t.Errorf("vcs should resolve to %s, got: %v", tb.branch, vcsinfo.Branch)
		}
		if err != nil {
			t.Error(err)
		}
	}
}

func TestFindsLocation(t *testing.T) {
	var configTests = []struct {
		url      string
		location string
	}{
		{"git+https://github.com/tony/.dot-configs@moo", "https://github.com/tony/.dot-configs"},
		{"git+ssh://git@github.com/tony/roundup.git@master", "ssh://git@github.com/tony/roundup.git"},
		{"hg+http://foicica.com/hg/textadept@ha", "http://foicica.com/hg/textadept"},
		{"svn+http://svn.code.sf.net/p/docutils/code/trunk@2019", "http://svn.code.sf.net/p/docutils/code/trunk"},
	}

	for _, tb := range configTests {
		vcsinfo, err := vcsync.ParsePIPUrl(tb.url)

		if vcsinfo.Location != tb.location {
			t.Errorf("vcs should resolve to %s, got: %v", tb.location, vcsinfo.Location)
		}
		if err != nil {
			t.Error(err)
		}
	}
}

func TestRepo(t *testing.T) {
	var configTests = []struct {
		url      string
		location string
		_type    vcs.Type
	}{
		{"git+https://github.com/tony/.dot-configs@moo", "https://github.com/tony/.dot-configs", vcs.Git},
		{"git+ssh://git@github.com/tony/roundup.git@master", "ssh://git@github.com/tony/roundup.git", vcs.Git},
		{"hg+http://foicica.com/hg/textadept@ha", "http://foicica.com/hg/textadept", vcs.Hg},
		{"svn+http://svn.code.sf.net/p/docutils/code/trunk@2019", "http://svn.code.sf.net/p/docutils/code/trunk", vcs.Svn},
	}

	tempDir, err := ioutil.TempDir("", "go-vcs-hg-tests")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
		}
	}()

	for _, tb := range configTests {
		vcsinfo, err := vcsync.ParsePIPUrl(tb.url)
		if err != nil {
			t.Error(err)
		}
		repo, err := vcsync.NewRepo(vcsinfo.VCSType, vcsinfo.Location, tempDir+"/testhgrepo")

		if err != nil {
			t.Error(err)
		}
		if repo.Vcs() != tb._type {
			t.Errorf("vcs should resolve to %s, got: %v", tb.location, vcsinfo.Location)
		}
	}
}
