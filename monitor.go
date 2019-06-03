package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/harrylee2015/monitor/conf"
	"github.com/harrylee2015/monitor/types"
	"github.com/harrylee2015/monitor/web"
	"os"
	"path/filepath"
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
	if _, err := toml.DecodeFile("monitor.toml", &cfg); err != nil {
		panic(err)
	}

	types.RegisterDB()

	server := web.NewWebServer(":8080")
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
