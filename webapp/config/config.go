package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/philippecery/libs/email"
)

type configStruct struct {
	Port              int
	Hostname          string
	TLS               *tlsConfigStruct
	DB                *dbConfigStruct
	Keys              *keysConfigStruct
	Email             *emailConfigStruct
	UserTokenValidity int
	DefaultLanguage   string
}

type dbConfigStruct struct {
	Name     string
	UserName string
	Password string
	Host     string
	Port     int
	TLS      *tlsConfigStruct
}

type tlsConfigStruct struct {
	Truststore string
	Keystore   string
	Password   string
}

type keysConfigStruct struct {
	UserID      string
	PII         string
	UserToken   string
	ActionToken string
}

type emailConfigStruct struct {
	Provider string
	Config   *email.ConfigStruct
	Bcc      string
}

var AppRoot string

// Config is loaded at application startup from provided JSON configuration file
var Config configStruct

func init() {
	if len(os.Args) < 2 {
		log.Fatal("Missing configuration file.")
	}
	var err error
	var configContent []byte
	if configContent, err = ioutil.ReadFile(os.Args[1]); err == nil {
		err = json.Unmarshal(configContent, &Config)
	}
	if err != nil {
		log.Fatalf("Unable to load configuration file. Cause: %v\n", err)
	}
}
