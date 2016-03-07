package main

import (
	"flag"
	"fmt"
	"reflect"

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
	dirs := viper.AllKeys()
	repos := []vcs.Repo{}
	for _, x := range dirs {

		v := viper.Sub(x)
		for _, k := range v.AllKeys() {
			fmt.Printf("%T\n", k)
			ExpandConfig(x, k, v.Get(k), &repos, v)
		}
	}
}

type verboseRepo struct {
	repo    string
	remotes map[string]interface{}
}

func castToMapStringInterface(
	src map[interface{}]interface{}) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[fmt.Sprintf("%v", k)] = v
	}
	return tgt
}

func ExpandConfig(dir string, name string, repoinfo interface{}, repos *[]vcs.Repo, v *viper.Viper) {
	fmt.Println(dir)

	switch repoinfo.(type) {
	case string:
	//
	case map[interface{}]interface{}:
		//
		fmt.Printf("%v\t: %v<%s>\n", name, repoinfo, reflect.TypeOf(repoinfo))
		vv := castToMapStringInterface(repoinfo.(map[interface{}]interface{}))
		// vv := v.Sub(name)
		fmt.Printf("created a sub sub: %v<%T>", vv, vv)
	case nil:
		return
	}
	if reflect.TypeOf(repoinfo).Kind() == reflect.Map {
		fmt.Printf("PMG ITS A MAP")
	}
}
