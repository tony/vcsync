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
	"strings"

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

// ParsePipURL parses PIP-style RFC3986 URLs.
func ParsePipURL(rawURL string) (VcsURL, error) {
	if strings.HasPrefix(rawURL, "git@github.com:") {
		rawURL = strings.Replace(rawURL, "git@github.com:", "git://github.com/", 1)
	}
	urlp, err := url.Parse(rawURL)
	if err != nil {
		return VcsURL{}, err
	}

	vre := regexp.MustCompile(`(?P<type>git|hg|svn|bzr)\+?(?P<location>.*)`)
	bre := regexp.MustCompile(`(?P<path>[^\@]*)@?(?P<branch>.*)$`)
	v := vre.FindStringSubmatch(urlp.Scheme)
	branch := bre.FindStringSubmatch(urlp.Path)

	var vtype vcs.Type
	// log.Infof("parsed url: %v, %v re:%+v", urlp.Path, urlp.Scheme, u)
	switch v[1] {
	case "git":
		vtype = vcs.Git
	case "hg":
		vtype = vcs.Hg
	case "svn":
		vtype = vcs.Svn
	case "default":
		return VcsURL{}, errors.New("Could not find VCS")
	}

	if v == nil {
		return VcsURL{}, ErrCannotDetectVCS
	}
	return VcsURL{
		*urlp, vtype, v[2] + "://" + urlp.Host + branch[1], branch[2],
	}, nil
}
