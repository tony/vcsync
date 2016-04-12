package vcsync_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Masterminds/vcs"
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
		repo, err := vcsync.NewRepo(vcsinfo.Vtype, vcsinfo.Location, tempDir+"/testhgrepo")

		if err != nil {
			t.Error(err)
		}
		if repo.Vcs() != tb.vtype {
			t.Errorf("vcs should resolve to %s, got: %v", tb.location, vcsinfo.Location)
		}
	}
}
