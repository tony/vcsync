// Support for parsing pip-style RFC3986 URL's
//
// Package provides support for pip (python package manager) style
// vcsync uses this style of URL for consolidating the VCS type,
// location (locally or internet) and branch in one string.

package vcsync

import (
	"errors"
	"net/url"
	"regexp"

	"github.com/tony/vcs"
)

// VcsURL stores parsed data from pip-style URLs.
type VcsURL struct {
	url.URL
	Vtype    vcs.Type
	Location string
	Branch   string
}

// Error for VCS detection and parsing failures
var (
	ErrCannotDetectVCS = errors.New("cannot detect VCS")
)

// ParsePIPUrl parses PIP-style RFC3986 URLs.
func ParsePIPUrl(rawURL string) (VcsURL, error) {
	urlp, err := url.Parse(rawURL)
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
	}

	return VcsURL{
		*urlp, vtype, u[2], u[3],
	}, nil
}
