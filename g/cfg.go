package g

import (
	"log"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/toolkits/file"
)

type Transfer struct {
	Enabled  bool
	Addrs    []string
	Interval int
	Timeout  int
}

type Http struct {
	Enabled  bool
	Listen   string
	Backdoor bool
}

type Global struct {
	Debug    bool
	Hostname string
	Addr     string
	Transfer *Transfer
	Http     *Http
}

var (
	ConfigFile string
	config     *Global
	lock   = new(sync.RWMutex)
)

func Config() *Global {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func Hostname() (string, error) {
	hostname := Config().Hostname
	if hostname != "" {
		return hostname, nil
	}
	if os.Getenv("FALCON_ENDPOINT") != "" {
		hostname = os.Getenv("FALCON_ENDPOINT")
		return hostname, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
	}
	return hostname, err
}

func IP() string {
	ip := Config().Addr
	if ip != "" {
		return ip
	}
	if len(LocalIp) > 0 {
		ip = LocalIp
	}
	return ip
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.toml cfg.toml`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var g  Global
	_, err = toml.Decode(configContent, &g)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()
	config = &g
	log.Println("read config file:", cfg, "successfully")
}
