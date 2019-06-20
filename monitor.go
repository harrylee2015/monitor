package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/harrylee2015/monitor/conf"
	"github.com/harrylee2015/monitor/types"
	"github.com/harrylee2015/monitor/web"
	"log"
	"net/http"
	"os"
	"path/filepath"
	_ "net/http/pprof"
	"sync"
)

var (
	configPath = flag.String("f", "monitor.toml", "configuration file")
	datadir    = flag.String("datadir", "", "data dir of monitor")
	group      sync.WaitGroup //定义一个同步等待的组
)

func main() {
	conf.CurrDir = pwd()
	var cfg conf.Config
	fmt.Println("dir:", conf.CurrDir)
	if _, err := toml.DecodeFile("monitor.toml", &cfg); err != nil {
		panic(err)
	}
	conf.SetConf(&cfg)
	types.RegisterDB()
	server := web.NewWebServer(":8080")
	go func() {
		log.Println(http.ListenAndServe(":10000",nil))
	}()
	server.Start()

	types.CloseDB()
}

func pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}
