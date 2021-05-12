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
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/util"
)

var Connection *sql.DB

const tlsName = "app2db"

func dataSource() string {
	dsnConfig := mysql.NewConfig()
	dsnConfig.User = config.Config.DB.UserName
	dsnConfig.Passwd = config.Config.DB.Password
	dsnConfig.Addr = fmt.Sprintf("%s:%s", config.Config.DB.Host, strconv.Itoa(config.Config.DB.Port))
	dsnConfig.DBName = config.Config.DB.Name
	dsnConfig.TLSConfig = tlsName
	dsnConfig.ParseTime = true
	dsnConfig.AllowAllFiles = false
	dsnConfig.AllowCleartextPasswords = false
	dsnConfig.AllowNativePasswords = true
	dsnConfig.AllowOldPasswords = false
	dsnConfig.MultiStatements = false
	return dsnConfig.FormatDSN()
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
					mysql.RegisterTLSConfig(tlsName, tlsConfig)
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
