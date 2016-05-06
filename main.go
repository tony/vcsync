package main

import (
	"fmt"
	"io"
	"os"
	"sync"

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

	log.Infof("%d repositories loaded.", len(r))

	tasks := make(chan *vcsync.VCSRepo, 64)

	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {

			for n := range tasks {
				cmd := vcsync.SyncRepo(n)
				stdout, err := cmd.StdoutPipe()
				stderr, err := cmd.StderrPipe()
				if err != nil {
					log.Fatal(err)
				}
				cmd.Start()
				go io.Copy(os.Stdout, stdout)
				go io.Copy(os.Stdout, stderr)
				cmd.Wait()
			}
			wg.Done()
		}()
	}

	// generate some tasks
	for _, m := range r[:10] {
		tasks <- &m
	}

	close(tasks)

	// wait for the workers to finish
	wg.Wait()
}
