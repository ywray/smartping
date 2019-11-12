package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"smartping/src/funcs"
	"smartping/src/g"
	"smartping/src/http"
	"time"

	"github.com/jakecoffman/cron"
	//"sync"
)

// Init config
var Version = "0.8.0"

// UpdateConfigFromFileFrequency 定时 load config.json 的频率（秒）
const UpdateConfigFromFileFrequency = 30

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	g.ParseConfig(Version)

	tick := time.Tick(time.Duration(UpdateConfigFromFileFrequency) * time.Second)
	go func() {
		for {
			select {
			case <-tick:
				g.UpDateConfigFromFile()
			}
		}
	}()

	go funcs.ClearArchive()
	c := cron.New()
	c.AddFunc("*/60 * * * * *", func() {
		go funcs.Ping()
		go funcs.Mapping()
		if g.Cfg.Mode["Type"] == "cloud" {
			go funcs.StartCloudMonitor()
		}
	}, "ping")
	c.AddFunc("0 0 * * * *", func() {
		go funcs.ClearArchive()
	}, "mtc")
	c.Start()
	http.StartHttp()
}
