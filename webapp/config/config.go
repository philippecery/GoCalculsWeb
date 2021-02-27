package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type configStruct struct {
	Port              int
	Hostname          string
	SSL               *sslConfigStruct
	DB                *dbConfigStruct
	Keys              *keysConfigStruct
	Gmail             *gmailConfigStruct
	UserTokenValidity int
	DefaultLanguage   string
}

type sslConfigStruct struct {
	Keystore string
	Password string
}

type dbConfigStruct struct {
	URL           string
	AuthSource    string
	AuthMechanism string
	UserName      string
	Password      string
}

type keysConfigStruct struct {
	UserID      string
	PII         string
	UserToken   string
	ActionToken string
}

type gmailConfigStruct struct {
	Oauth2 *oauth2ConfigStruct
	Bcc    string
}

type oauth2ConfigStruct struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
}

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
