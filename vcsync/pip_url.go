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
	VCS string
}

var (
	ErrCannotDetectVCS = errors.New("Cannot detect VCS")
)

func ParsePIPUrl(vcsUrl string) (VcsURL, error) {
	urlp, err := url.Parse(vcsUrl)
	if err != nil {
		return VcsURL{}, err
	}
	v := regexp.MustCompile(`(?P<type>git|hg|svn|bzr).*$`)

	m := v.FindStringSubmatch(urlp.Scheme)
	if m == nil {
		return VcsURL{}, ErrCannotDetectVCS
	} else {
		return VcsURL{
			*urlp, m[1],
		}, nil
	}
}
