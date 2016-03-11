package vcsync_test

import "testing"
import "github.com/tony/vcsync/vcsync"

func TestFindsGit(t *testing.T) {
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

		vcstype, err := vcsync.ParsePIPUrl(tt.url)

		if vcstype.VCS != tt.vcs {
			t.Errorf("vcs should resolve to git, got: %v", vcstype)
		}
		if err != nil {
			t.Error(err)
		}
	}
}
