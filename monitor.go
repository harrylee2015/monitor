package main

import (
	"flag"
	"os"
	"path/filepath"
	"sync"
    "github.com/BurntSushi/toml"
	"github.com/harrylee2015/monitor/model"
	"github.com/harrylee2015/monitor/web"
)

var (
	configPath = flag.String("f", "monitor.toml", "configuration file")
	datadir    = flag.String("datadir", "", "data dir of monitor")
	group      sync.WaitGroup //定义一个同步等待的组
)

func main() {
	var cfg model.Config
	if _,err := toml.DecodeFile("monitor.toml",&cfg);err !=nil{
		panic(err)
	}

	model.RegisterDB()

	server :=web.NewWebServer(":8080")
	server.Start()

	model.CloseDB()
}

func pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}
