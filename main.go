package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tony/vcsync/vcsync"
)

func main() {
	viper.SetConfigName(".vcspull")
	viper.AddConfigPath("$HOME")
	var r vcsync.Repos
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	m := map[string]map[string]interface{}{}

	for _, x := range viper.AllKeys() {
		m[x] = viper.GetStringMap(x)
		r.LoadRepos(x, m[x])
		log.Debug(r[len(r)-1:])
	}

	for _, repo := range r {
		log.Infof("%s @ %s", repo.Repo.LocalPath(), repo.Repo.Remote())
	}
	log.Infof("%d repositories loaded.", len(r))
}
