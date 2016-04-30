package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tony/vcsync/vcsync"
)

var r vcsync.Repos

func main() {
	viper.SetConfigName(".vcspull")
	viper.AddConfigPath("$HOME")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
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
