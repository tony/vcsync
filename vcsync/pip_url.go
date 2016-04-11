// Package provides support for pip (python package manager) style
// URL's

package vcsync

import (
	"errors"
	"net/url"
	"regexp"
)

type VcsURL struct {
	url.URL
	VCSType  string
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

	if u == nil {
		return VcsURL{}, ErrCannotDetectVCS
	} else {
		return VcsURL{
			*urlp, u[1], u[2], u[3],
		}, nil
	}
}
