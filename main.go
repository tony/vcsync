package main

import (
	"fmt"
	"net/url"

	"github.com/Masterminds/vcs"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tony/vcsync/vcsync"
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
	var legacyRepos []vcsync.LegacyRepoConf
	for _, x := range viper.AllKeys() {
		m[x] = viper.GetStringMap(x)
		vcsync.ExpandConfig(x, m[x], &legacyRepos)
		log.Debug(legacyRepos[len(legacyRepos)-1:])
	}
	log.Infof("%d repositories loaded.", len(legacyRepos))

	for _, repo := range legacyRepos {
		repoURL, err := url.Parse(repo.URL)
		if err != nil {
			log.Infof("Error parsing URL %v", repo.URL)
		}
		repoURL.Scheme = ""
		log.Info(repoURL.String())
		log.Info(repo.URL)
	}
}
