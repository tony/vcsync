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
		ExpandConfig(x, m[x], &repos)
	}
}

type verboseRepo struct {
	repo    string
	remotes map[string]interface{}
}

func ExpandConfig(dir string, entries map[string]interface{}, repos *[]vcs.Repo) {
	fmt.Println(dir)
	for name, repo := range entries {
		switch repo.(type) {
		case string:
			fmt.Printf("name: %v\t repo: %v\n", name, repo)
		case map[string]verboseRepo:
			fmt.Printf("string nested name: %v\t repo: %v (%T)\n", name, repo, repo)
		case map[string]interface{}:
			fmt.Printf("string nested name: %v\t repo: %v (%T)\n", name, repo, repo)
		case map[interface{}]interface{}:
			repo = castToMapStringInterface(repo.(map[interface{}]interface{}))
			fmt.Println(repo)
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
