package g

import (
	"log"
	"os"
	// "github.com/flystary/agent/model"
)

var Root string

func InitRootDir() {
	var err error
	Root, err = os.Getwd()
	if err != nil {
		log.Fatalln("getwd fail:", err)
	}
}

var LocalIp string

func IsTrustable(ip string) bool {
	return true
}
