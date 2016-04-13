package vcsync

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/tony/vcs"
)

type LegacyRepoConf struct {
	Name    string
	Url     string
	Path    string
	Remotes map[string]string
}

// Version retrieves the current version.
func AddRemote(s *vcs.GitRepo, name, url string) (string, error) {
	out, err := s.RunFromDir("git", "remote", "add", name, url)

	if err != nil {
		if strings.Contains(string(out), fmt.Sprintf("remote %s already exists.", name)) {
			return "", errors.New(string(out))
			// 		return UpdateRemote(s, name, url)
		}
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func UpdateRemote(s *vcs.GitRepo, name, url string) (string, error) {
	out, err := s.RunFromDir("git", "remote", "set-url", name, url)
	if err != nil {
		// if string(out) != "" {
		// 	return "", errors.New(strings.TrimSpace(string(out)))
		// }
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func ExpandConfig(dir string, entries map[string]interface{}, repos *[]LegacyRepoConf) {
	for name, repo := range entries {
		log.Debug("name: %v\t repo: %v", name, repo)
		legacyRepo := LegacyRepoConf{
			Name: name,
			Path: dir,
		}
		switch repo.(type) {
		case string:
			legacyRepo.Url = repo.(string)
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
			legacyRepo.Url = r["repo"].(string)

		default:
			log.Infof("undefined name %v: verbose repo (type %T)\n", name, repo)
			continue
		}
		*repos = append(*repos, legacyRepo)
	}
}

func NewRepo(vtype vcs.Type, remote, local string) (vcs.Repo, error) {
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

func GetRepo(repo vcs.Repo) {
	repo.Vcs()
	// Returns Git as this is a Git repo

	err := repo.Get()
	// Pulls down a repo, or a checkout in the case of SVN, and returns an
	// error if that didn't happen successfully.
	if err != nil {
		fmt.Println(err)
	}

	err = repo.UpdateVersion("master")
	// Checkouts out a specific version. In most cases this can be a commit id,
	// branch, or tag.
	if err != nil {
		fmt.Println(err)
	}
}
