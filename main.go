package main

import (
	"fmt"

	"github.com/Masterminds/vcs"
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
	fmt.Println(repos)
	var legacyRepos []LegacyRepoConf
	for _, x := range viper.AllKeys() {
		m[x] = viper.GetStringMap(x)
		ExpandConfig(x, m[x], &legacyRepos, viper.Sub(x))
		fmt.Println(legacyRepos[len(legacyRepos)-1:][0])
	}
}

type LegacyRepoConf struct {
	name    string
	url     string
	path    string
	remotes map[string]string
}

func ExpandConfig(dir string, entries map[string]interface{}, repos *[]LegacyRepoConf, v *viper.Viper) {
	//fmt.Println(dir)
	for name, repo := range entries {
		// fmt.Printf("name: %v\t repo: %v\n", name, repo)
		var remotes map[string]string
		var repo_url string
		switch repo.(type) {
		case string:
			remotes = nil
			repo_url = repo.(string)
		case map[interface{}]interface{}:
			r := cast.ToStringMap(repo.(map[interface{}]interface{}))
			remote_map := cast.ToStringMapString(r["remotes"].(map[interface{}]interface{}))
			if r["remotes"] != nil {
				for remote_name, remote := range remote_map {
					remotes = map[string]string{
						remote_name: remote,
					}
				}
			} else {
				remotes = nil
				fmt.Printf("No remotes detected, check your formatting for %s at %s", name, repo)
			}
			repo_url = r["repo"].(string)

		default:
			fmt.Printf("undefined name %v: verbose repo (type %T)\n", name, repo)
			continue
		}
		*repos = append(*repos, LegacyRepoConf{
			name:    name,
			path:    dir,
			url:     repo_url,
			remotes: remotes,
		})
	}
}
