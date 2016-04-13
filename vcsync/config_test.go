package vcsync_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/tony/vcs"
	"github.com/tony/vcsync/vcsync"
)

func TestRepo(t *testing.T) {
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

	for _, tb := range configTests {
		vcsinfo, err := vcsync.ParsePIPUrl(tb.url)
		if err != nil {
			t.Error(err)
		}
		repo, err := vcsync.NewRepo(vcsinfo.Vtype, vcsinfo.Location, tempDir+"/testhgrepo")

		if err != nil {
			t.Error(err)
		}
		if repo.Vcs() != tb.vtype {
			t.Errorf("vcs should resolve to %s, got: %v", tb.location, vcsinfo.Location)
		}
	}
}

func TestGitRemotes(t *testing.T) {
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
	remotePath := tempDir + "/VCSTestRemote"
	repoPath := tempDir + "/VCSTestRepo"

	// init repo
	// initCmd := exec.Command("git", "init", repoPath)
	// _, err = initCmd.CombinedOutput()

	rmtInitCmd := exec.Command("git", "init", remotePath)
	_, err = rmtInitCmd.CombinedOutput()

	if err != nil {
		t.Error(err)
	}

	repo, err := vcs.NewGitRepo(remotePath, repoPath)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Init()
	if err != nil {
		t.Fatal(err)
	}

	new_repo_remote := "https://github.com/lol/lol"

	_, err = vcsync.AddRemote(repo, "origin", new_repo_remote)
	if err != nil {
		t.Error(err)
	}

	out, err := repo.RunFromDir("git", "remote", "-v")
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(out), new_repo_remote) {
		t.Errorf("vcs should update remote properly, %s not found in %s", new_repo_remote, string(out))
	}
}
