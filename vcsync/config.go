package vcsync

import (
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cast"
)

// Repos is a Slice of Repo(s)
type Repos []VCSRepo

// LoadRepos expands the JSON/YAML configuration into Repo objects.
func (r *Repos) LoadRepos(dir string, items map[string]interface{}) {
	var err error
	for name, data := range items {
		var repoURL string
		log.Debug("name: %v\t data: %v", name, data)
		repo := VCSRepo{
			Name: name,
		}
		switch data.(type) {
		case string:
			repoURL = data.(string)
		case map[interface{}]interface{}:
			r := cast.ToStringMap(data)
			if r["remotes"] != nil {
				repo.Remotes = make(map[string]string)
				for rname, rurl := range cast.ToStringMapString(r["remotes"]) {
					repo.Remotes[rname] = rurl
				}
			} else {
				log.Infof("No remotes detected, check your formatting for %s at %s", name, data)
			}
			repoURL = r["repo"].(string)

		default:
			log.Infof("undefined name %v: verbose repo (type %T)\n", name, data)
			continue
		}

		repo.Repo, err = NewRepoFromPipURL(repoURL, path.Join(dir, name))
		if err != nil {
			log.Infof("failure adding %v (type %T) as repo\n", name, data)
			continue
		}
		*r = append(*r, repo)
	}
}
