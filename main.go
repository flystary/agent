package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/flystary/agent/g"
	"github.com/flystary/agent/http"
	log "github.com/sirupsen/logrus"
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

	if g.SystemType != "windows" {
		watch()
	}

	g.InitRootDir()

	go http.Start()

	select {}
}

func watch() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL) // 其中 SIGKILL = kill -9 <pid> 可能无法截获
	go func() {
		log.Debug("watching stop signals")
		<-c
		log.Debug("ready to exit on SIGTERM")
		os.Exit(1)
	}()
}
