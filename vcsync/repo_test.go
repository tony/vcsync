package vcsync_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/Masterminds/vcs"
	"github.com/tony/vcsync/vcsync"
)

func TestNewRepo(t *testing.T) {
	var configTests = []struct {
		url      string
		location string
		vtype    vcs.Type
	}{
		{"git+https://github.com/tony/.dot-configs@moo", "https://github.com/tony/.dot-configs", vcs.Git},
		{"git+ssh://git@github.com/tony/roundup.git@master", "ssh://git@github.com/tony/roundup.git", vcs.Git},
		{"hg+http://foicica.com/hg/textadept@ha", "http://foicica.com/hg/textadept", vcs.Hg},
		{"svn+http://svn.code.sf.net/p/docutils/code/trunk@2019", "http://svn.code.sf.net/p/docutils/code/trunk", vcs.Svn},
	}

	for _, tb := range configTests {
		vURL, err := vcsync.ParsePipURL(tb.url)
		if err != nil {
			t.Error(err)
		}
		repo, err := vcsync.NewRepo(vURL.Vtype, vURL.Location(), "./testhgrepo")

		if err != nil {
			t.Error(err)
		}
		if repo.Vcs() != tb.vtype {
			t.Errorf("vcs should resolve to %s, got: %v", tb.location, vURL.Location())
		}
	}
}

func createTempGitRepo(localPath, remotePath string) (*vcs.GitRepo, error) {
	repo, err := vcs.NewGitRepo(localPath, remotePath)
	if err != nil {
		return nil, err
	}

	if err = repo.Init(); err != nil {
		return nil, err
	}

	return repo, nil
}

func TestGitAddRemote(t *testing.T) {
	new_repo_remote := "https://github.com/lol/lol"
	second_remote := "https://github.com/wut/wut"
	tempDir, err := ioutil.TempDir("", "go-vcs-tests")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
		}
	}()

	repo, err := createTempGitRepo(tempDir+"/VCSTestRemote", tempDir+"/VCSTestRepo")
	if err != nil {
		t.Error(err)
	}

	if _, err = vcsync.AddRemote(repo, "origin", new_repo_remote); err != nil {
		t.Error(err)
	}

	out, err := repo.RunFromDir("git", "remote", "-v")
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(out), new_repo_remote) {
		t.Errorf("vcs should add remote properly, %s not found in %s", new_repo_remote, string(out))
	}

	// Adding a remote already exists
	if _, err = vcsync.AddRemote(repo, "origin", second_remote); err != nil {
		if !strings.Contains(err.Error(), fmt.Sprintf("remote %s already exists.", "origin")) {
			t.Error("Adding a remote if one exist should raise an error")
		}
	}
}

func TestGitUpdateRemote(t *testing.T) {
	second_remote := "https://github.com/wut/wut"
	tempDir, err := ioutil.TempDir("", "go-vcs-tests")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
		}
	}()

	repo, err := createTempGitRepo(tempDir+"/VCSTestRemote", tempDir+"/VCSTestRepo")
	if err != nil {
		t.Error(err)
	}

	new_repo_remote := "https://github.com/lol/lol"

	if _, err = vcsync.AddRemote(repo, "origin", new_repo_remote); err != nil {
		t.Error(err)
	}

	// Update a repo
	if _, err = vcsync.UpdateRemote(repo, "origin", second_remote); err != nil {
		t.Error(err)
	}

	out, err := repo.RunFromDir("git", "remote", "-v")
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(out), second_remote) {
		t.Errorf("vcs should update remote properly, %s not found in %s", second_remote, string(out))
	}
}

func TestSyncRepo(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "go-vcs-tests")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
		}
	}()

	var configTests = []struct {
		localPath  string
		remotePath string
		vtype      vcs.Type
	}{
		{tempDir + "/localGit", tempDir + "/remoteGit", vcs.Git},
	}

	repos := make([]vcs.Repo, 0)
	for _, c := range configTests {
		r, err := vcsync.NewRepo(c.vtype, c.localPath, c.remotePath)
		if err != nil {
			t.Error(err)
		}
		if _, err := createTempGitRepo(c.localPath, c.remotePath); err != nil {
			t.Error(err)
		}
		repos = append(repos, r)
	}

	for _, r := range repos {
		// creates repo if doesn't exist
		cmd := vcsync.SyncRepo(r)
		if r.GetCmd() == cmd {
			t.Errorf("want %s, got %s", r.GetCmd(), cmd)
		}

		_, err = os.Stat(r.LocalPath())
		if err != nil {
			t.Error(err)
		}

		vcsync.SyncRepo(r)
		// updates repo if exists
	}
}
