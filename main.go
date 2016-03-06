package main

import (
	"flag"
	"fmt"

	"github.com/Masterminds/vcs"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	viper.SetConfigName(".vcspull")
	viper.AddConfigPath("$HOME")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	m := map[string]map[string]interface{}{}
	dirs := viper.AllKeys()
	repos := []vcs.Repo{}
	for _, x := range dirs {
		m[x] = viper.GetStringMap(x)
		ExpandConfig(x, m[x], &repos)
	}
}

func ExpandConfig(dir string, repositories map[string]interface{}, repos *[]vcs.Repo) {
	fmt.Println(dir)
	for name, repo := range repositories {
		fmt.Println(name, repo)
	}
}
