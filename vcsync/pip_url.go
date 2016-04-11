// Package provides support for pip (python package manager) style
// URL's

package vcsync

import (
	"errors"
	"net/url"
	"regexp"

	"github.com/Masterminds/vcs"
)

type VcsURL struct {
	url.URL
	VCSType  vcs.Type
	Location string
	Branch   string
}

var (
	ErrCannotDetectVCS = errors.New("Cannot detect VCS")
)

func ParsePIPUrl(vcsUrl string) (VcsURL, error) {
	urlp, err := url.Parse(vcsUrl)
	if err != nil {
		return VcsURL{}, err
	}

	v := regexp.MustCompile(`(?P<type>git|hg|svn|bzr)\+(?P<location>.*?)@?(?P<branch>[^@]*)$`)
	u := v.FindStringSubmatch(urlp.String())

	var vtype vcs.Type
	switch u[1] {
	case "git":
		vtype = vcs.Git
	case "hg":
		vtype = vcs.Hg
	case "svn":
		vtype = vcs.Svn
	}

	if u == nil {
		return VcsURL{}, ErrCannotDetectVCS
	} else {
		return VcsURL{
			*urlp, vtype, u[2], u[3],
		}, nil
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
