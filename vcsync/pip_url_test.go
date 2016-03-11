package vcsync_test

import "testing"
import "github.com/tony/vcsync/vcsync"

func TestFindsGit(t *testing.T) {
	vcstype, err := vcsync.ParsePIPUrl("git+https://github.com/tony/.dot-configs")

	if vcstype.VCS != "git" {
		t.Errorf("vcs should resolve to git, got: %v", vcstype)
	}
	if err != nil {
		t.Error(err)
	}
	// Output:
	// Hi
}
