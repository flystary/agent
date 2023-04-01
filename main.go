package main

import (
	"github.com/flystary/agent/g"
	"flag"
	"fmt"
	"os"
	"github.com/flystary/agent/http"
)

func main() {
	//cfg
	cfg := flag.String("c", "cfg.toml", "configuretion file")
	version := flag.Bool("v", false, "agentplay Version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	if g.Config().Debug {
		g.InitLog("debug")
	} else {
		g.InitLog("info")
	}

	g.InitRootDir()
	go http.Start()

	select {}
}
