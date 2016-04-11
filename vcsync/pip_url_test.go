package vcsync_test

import "testing"
import "github.com/tony/vcsync/vcsync"

func TestFindsVcsType(t *testing.T) {
	var configTests = []struct {
		url string
		vcs string
	}{
		{"git+https://github.com/tony/.dot-configs", "git"},
		{"git+ssh://git@github.com/tony/roundup.git", "git"},
		{"hg+http://foicica.com/hg/textadept", "hg"},
		{"svn+http://svn.code.sf.net/p/docutils/code/trunk", "svn"},
	}

	for _, tt := range configTests {
		vcsinfo, err := vcsync.ParsePIPUrl(tt.url)

		if vcsinfo.VCSType != tt.vcs {
			t.Errorf("vcs should resolve to %s, got: %v", tt.vcs, vcsinfo.VCSType)
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
