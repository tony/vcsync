package main

import (
	"fmt"

	"github.com/Masterminds/vcs"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName(".vcspull")
	viper.AddConfigPath("$HOME")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	m := map[string]map[string]interface{}{}
	var repos []vcs.Repo
	log.Println(repos)
	var legacyRepos []LegacyRepoConf
	for _, x := range viper.AllKeys() {
		m[x] = viper.GetStringMap(x)
		ExpandConfig(x, m[x], &legacyRepos, viper.Sub(x))
		log.Println(legacyRepos[len(legacyRepos)-1:])
	}
	log.Printf("%d repositories loaded.", len(legacyRepos))
}

type LegacyRepoConf struct {
	name    string
	url     string
	path    string
	remotes map[string]string
}

func ExpandConfig(dir string, entries map[string]interface{}, repos *[]LegacyRepoConf, v *viper.Viper) {
	for name, repo := range entries {
		log.Debug("name: %v\t repo: %v", name, repo)
		legacyRepo := LegacyRepoConf{
			name: name,
			path: dir,
		}
		switch repo.(type) {
		case string:
			legacyRepo.url = repo.(string)
		case map[interface{}]interface{}:
			r := cast.ToStringMap(repo)
			if r["remotes"] != nil {
				legacyRepo.remotes = make(map[string]string)
				for rname, rurl := range cast.ToStringMapString(r["remotes"]) {
					legacyRepo.remotes[rname] = rurl
				}
			} else {
				log.Printf("No remotes detected, check your formatting for %s at %s", name, repo)
			}
			legacyRepo.url = r["repo"].(string)

		default:
			log.Printf("undefined name %v: verbose repo (type %T)\n", name, repo)
			continue
		}
		*repos = append(*repos, legacyRepo)
	}
}
