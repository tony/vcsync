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
		switch repo.(type) {
		case string:
			*repos = append(*repos, LegacyRepoConf{
				name:    name,
				url:     repo.(string),
				path:    dir,
				remotes: nil,
			})
		case map[interface{}]interface{}:
			r := castToMapStringInterface(repo.(map[interface{}]interface{}))
			if r["remotes"] != nil {
				for remote_name, remote := range castToMapStringString(r["remotes"].(map[interface{}]interface{})) {
					*repos = append(*repos, LegacyRepoConf{
						name: name,
						path: dir,
						url:  r["repo"].(string),
						remotes: map[string]string{
							remote_name: remote,
						},
					})
				}
			} else {
				*repos = append(*repos, LegacyRepoConf{
					name:    name,
					path:    dir,
					url:     r["repo"].(string),
					remotes: nil,
				})
				fmt.Printf("No remotes detected, check your formatting for %s at %s", name, repo)
			}
		default:
			fmt.Printf("name %v: verbose repo (type %T)\n", name, repo)
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

func castToMapStringString(src map[interface{}]interface{}) map[string]string {
	tgt := make(map[string]string)
	for k, v := range src {
		tgt[fmt.Sprintf("%v", k)] = v.(string)
	}
	return tgt
}
