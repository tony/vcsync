// Support for parsing pip-style RFC3986 URL's
//
// Package provides support for pip (python package manager) style
// vcsync uses this style of URL for consolidating the VCS type,
// location (locally or internet) and ref in one string.

package vcsync

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/Masterminds/vcs"
	log "github.com/Sirupsen/logrus"
)

// VcsURL stores parsed data from pip-style URLs. It is built on top of url.URL in stdlib.
type VcsURL struct {
	*url.URL
	Vtype vcs.Type
	Ref   string
}

// ParsePipURL parses PIP-style RFC3986 URLs.
func ParsePipURL(rawURL string) (VcsURL, error) {
	if strings.HasPrefix(rawURL, "git@github.com:") {
		rawURL = strings.Replace(rawURL, "git@github.com:", "git://github.com/", 1)
	}
	urlp, err := url.Parse(rawURL)

	if err != nil {
		return VcsURL{}, err
	}

	vcsURL := VcsURL{
		URL: urlp,
	}

	// For finding the VCS type and scheme + splitting them
	// git+https
	// git, https
	vre := regexp.MustCompile(`(?P<type>git|hg|svn|bzr)\+?(?P<location>.*)`)

	// For grabbing the refs in URL's, if they exist, e.g.
	// git+https://github.com/tony/vcsync@develop
	bre := regexp.MustCompile(`(?P<path>[^\@]*)@?(?P<ref>.*)$`)

	v := vre.FindStringSubmatch(vcsURL.URL.Scheme)
	ref := bre.FindStringSubmatch(vcsURL.URL.Path)

	if v == nil {
		return VcsURL{}, vcs.ErrCannotDetectVCS
	}

	switch v[1] {
	case "git":
		vcsURL.Vtype = vcs.Git
	case "hg":
		vcsURL.Vtype = vcs.Hg
	case "svn":
		vcsURL.Vtype = vcs.Svn
	case "default":
		return VcsURL{}, vcs.ErrCannotDetectVCS
	}

	vcsURL.URL.Path = ref[1]
	vcsURL.URL.Scheme = v[2]
	vcsURL.Ref = ref[2]
	log.Infof("vcsurl: %+v", vcsURL)

	return vcsURL, nil
}

// Location returns the location of the repository.
func (v *VcsURL) Location() string {
	return v.Scheme + "://" + v.Host + v.Path
}
