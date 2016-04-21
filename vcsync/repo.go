package vcsync

// VCSRepo is an ephemeral data struct for processing configs.
import (
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/vcs"
	log "github.com/Sirupsen/logrus"
)

// VCSRepo holds repo name, Repo object and list of remotes
type VCSRepo struct {
	vcs.Repo
	Name    string
	Remotes map[string]string
}

// AddRemote adds a remote to a git repository.
func AddRemote(s *vcs.GitRepo, name, url string) (string, error) {
	out, err := s.RunFromDir("git", "remote", "add", name, url)

	if err != nil {
		if strings.Contains(string(out), fmt.Sprintf("remote %s already exists.", name)) {
			return "", errors.New(string(out))
		}
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// UpdateRemote updates the url of a current remote.
func UpdateRemote(s *vcs.GitRepo, name, url string) (string, error) {
	out, err := s.RunFromDir("git", "remote", "set-url", name, url)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// NewRepo is a generic function for created a new repo object from vcs.Type.
func NewRepo(vtype vcs.Type, remote, local string) (vcs.Repo, error) {
	log.Info(remote)
	switch vtype {
	case vcs.Git:
		return vcs.NewGitRepo(remote, local)
	case vcs.Svn:
		return vcs.NewSvnRepo(remote, local)
	case vcs.Hg:
		return vcs.NewHgRepo(remote, local)
	case vcs.Bzr:
		return vcs.NewBzrRepo(remote, local)
	}

	// Should never fall through to here but just in case.
	return nil, vcs.ErrCannotDetectVCS
}

// NewRepoFromPipURL returns Repo object from pip url
func NewRepoFromPipURL(remote, local string) (vcs.Repo, error) {
	pipURL, err := ParsePipURL(remote)
	log.Infof("%+v", pipURL)
	if err != nil {
		return nil, err
	}

	return NewRepo(pipURL.Vtype, pipURL.Location(), local)
}
