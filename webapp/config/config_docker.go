// +build docker

package config

import (
	"log"
)

func init() {
	log.Println("Docker installation")
	AppRoot = "/"
}
