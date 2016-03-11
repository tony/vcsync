package vcsync

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cast"
)

type LegacyRepoConf struct {
	Name    string
	Url     string
	Path    string
	Remotes map[string]string
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
