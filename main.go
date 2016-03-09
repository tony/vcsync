package main

import (
	"fmt"

	"github.com/Masterminds/vcs"
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
	repos := []vcs.Repo{}
	for _, x := range viper.AllKeys() {
		m[x] = viper.GetStringMap(x)
		ExpandConfig(x, m[x], &repos, viper.Sub(x))
	}
}

type RepoConf struct {
	name    string
	url     string
	path    string
	remotes map[string]string
}

func ExpandConfig(dir string, entries map[string]interface{}, repos *[]vcs.Repo, v *viper.Viper) {
	//fmt.Println(dir)
	for name, repo := range entries {
		// fmt.Printf("name: %v\t repo: %v\n", name, repo)
		switch repo.(type) {
		case string:
			repoConf := RepoConf{
				name:    name,
				url:     repo.(string),
				path:    dir,
				remotes: nil,
			}
			fmt.Println(repoConf)
		case map[interface{}]interface{}:
			r := castToMapStringInterface(repo.(map[interface{}]interface{}))
			if r["remotes"] != nil {
				for remote_name, remote := range r["remotes"].(map[interface{}]interface{}) {
					fmt.Println(RepoConf{
						name: name,
						path: dir,
						url:  r["repo"].(string),
						remotes: map[string]string{
							remote_name.(string): remote.(string),
						},
					})
				}
			}
		default:
			// fmt.Printf("name %v: verbose repo (type %T)\n", name, repo)
		}
	}
}

func castToMapStringInterface(src map[interface{}]interface{}) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[fmt.Sprintf("%v", k)] = v
	}
	return tgt
}
