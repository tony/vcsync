package vcsync_test

import (
	"testing"

	"github.com/Masterminds/vcs"
	"github.com/tony/vcsync/vcsync"
)

func TestLoadRepo(t *testing.T) {
	var r vcsync.Repos
	r.LoadRepos("dir", map[string]interface{}{"myrepo": "git+https://github.com/moo/moo"})

	if len(r) != 1 {
		t.Errorf("loaded %v instead of one", len(r))
	}
	if r[0].Repo.Vcs() != vcs.Git {
		t.Errorf("should return correct repo vcs type, found %v", r[0].Repo.Vcs())
	}
}

func TestLoadExpanded(t *testing.T) {
	var r vcsync.Repos
	r.LoadRepos(
		"dir",
		map[string]interface{}{"myrepo": map[interface{}]interface{}{
			"repo": "hg+https://github.com/moo/moo",
		},
		},
	)

	if len(r) != 1 {
		t.Errorf("loaded %v instead of one", len(r))
	}
	if r[0].Repo.Vcs() != vcs.Hg {
		t.Errorf("should return correct repo vcs type, found %v", r[0].Repo.Vcs())
	}
}
