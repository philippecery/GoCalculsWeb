package database

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/util"
)

var Connection *sql.DB

var parametersMap = map[string]string{
	"tls":                     "app2db",
	"allowAllFiles":           "false",
	"allowCleartextPasswords": "false",
	"allowNativePasswords":    "false",
	"allowOldPasswords":       "false",
	"multiStatements":         "false",
}

func dataSource() string {
	parameters := []string{}
	for k, v := range parametersMap {
		parameters = append(parameters, k+"="+v)
	}
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?%s", config.Config.DB.UserName, config.Config.DB.Password, config.Config.DB.Host, strconv.Itoa(config.Config.DB.Port), config.Config.DB.Name, strings.Join(parameters, "&"))
}

func Open() error {
	var err error
	if config.Config.DB.TLS.Truststore != "" {
		var pem []byte
		rootCertPool := x509.NewCertPool()
		if pem, err = ioutil.ReadFile(config.Config.DB.TLS.Truststore); err == nil {
			if rootCertPool.AppendCertsFromPEM(pem) {
				tlsConfig := &tls.Config{
					RootCAs: rootCertPool,
				}
				if config.Config.DB.TLS.Keystore != "" {
					if config.Config.DB.TLS.Password != "" {
						var clientKeyPass = []byte(config.Config.DB.TLS.Password)
						var keyPair tls.Certificate
						if keyPair, err = util.X509KeyPair(config.Config.DB.TLS.Keystore, &clientKeyPass); err == nil {
							tlsConfig.Certificates = []tls.Certificate{keyPair}
						}
					} else {
						err = errors.New("database: missing private key password")
					}
				}
				if err == nil {
					mysql.RegisterTLSConfig(parametersMap["tls"], tlsConfig)
					if Connection, err = sql.Open("mysql", dataSource()); err == nil {
						if err = Connection.Ping(); err == nil {
							Connection.SetConnMaxLifetime(time.Minute * 3)
							Connection.SetMaxOpenConns(10)
							Connection.SetMaxIdleConns(10)
							log.Printf("database: connected")
						}
					}
				}
			}
		}
	} else {
		err = errors.New("database: missing TLS truststore")
	}
	return err
}

func Close() error {
	return Connection.Close()
}
