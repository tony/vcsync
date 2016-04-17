package vcsync

import (
	"errors"
	"fmt"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/tony/vcs"
)

// VCSRepo is an ephemeral data struct for processing configs.
type VCSRepo struct {
	Repo    vcs.Repo
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

// ExpandConfig expands the JSON/YAML configuration into Repo objects.
func ExpandConfig(dir string, entries map[string]interface{}, repos *[]VCSRepo) {
	for name, repo := range entries {
		var repoURL string
		log.Debug("name: %v\t repo: %v", name, repo)
		legacyRepo := VCSRepo{
			Name: name,
		}
		switch repo.(type) {
		case string:
			repoURL = repo.(string)
		case map[interface{}]interface{}:
			r := cast.ToStringMap(repo)
			if r["remotes"] != nil {
				legacyRepo.Remotes = make(map[string]string)
				for rname, rurl := range cast.ToStringMapString(r["remotes"]) {
					legacyRepo.Remotes[rname] = rurl
				}
			} else {
				log.Infof("No remotes detected, check your formatting for %s at %s", name, repo)
			}
			repoURL = r["repo"].(string)

		default:
			log.Infof("undefined name %v: verbose repo (type %T)\n", name, repo)
			continue
		}
		var err error
		legacyRepo.Repo, _ = NewRepoFromPipURL(repoURL, path.Join(dir, name))
		if err != nil {
			log.Infof("failure adding %v (type %T) as repo\n", name, repo)
			continue
		}
		*repos = append(*repos, legacyRepo)
	}
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
	return nil, ErrCannotDetectVCS
}
