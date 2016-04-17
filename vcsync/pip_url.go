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

	log "github.com/Sirupsen/logrus"
	"github.com/tony/vcs"
)

// VcsURL stores parsed data from pip-style URLs.
type VcsURL struct {
	*url.URL
	Vtype  vcs.Type
	Branch string
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

	vcsURL := VcsURL{
		URL: urlp,
	}

	// For finding the VCS type and scheme + splitting them
	// git+https
	// git, https
	vre := regexp.MustCompile(`(?P<type>git|hg|svn|bzr)\+?(?P<location>.*)`)

	// For grabbing the refs in URL's, if they exist, e.g.
	// git+https://github.com/tony/vcsync@develop
	bre := regexp.MustCompile(`(?P<path>[^\@]*)@?(?P<branch>.*)$`)

	v := vre.FindStringSubmatch(vcsURL.URL.Scheme)
	branch := bre.FindStringSubmatch(vcsURL.URL.Path)

	if v == nil {
		return VcsURL{}, ErrCannotDetectVCS
	}

	switch v[1] {
	case "git":
		vcsURL.Vtype = vcs.Git
	case "hg":
		vcsURL.Vtype = vcs.Hg
	case "svn":
		vcsURL.Vtype = vcs.Svn
	case "default":
		return VcsURL{}, ErrCannotDetectVCS
	}

	vcsURL.URL.Path = branch[1]
	vcsURL.URL.Scheme = v[2]
	vcsURL.Branch = branch[2]
	log.Infof("vcsurl: %+v", vcsURL)

	return vcsURL, nil
}

func (v *VcsURL) Location() string {
	return v.Scheme + "://" + v.Host + v.Path
}
