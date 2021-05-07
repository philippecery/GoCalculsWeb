package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/util"
)

func create() (*http.Server, error) {
	if config.Config.TLS == nil {
		return nil, errors.New("server: missing TLS configuration")
	}
	if config.Config.TLS.Password == "" {
		return nil, errors.New("server: missing private key password")
	}
	var serverKeyPass = []byte(config.Config.TLS.Password)
	keyPair, err := util.X509KeyPair(config.Config.TLS.Keystore, &serverKeyPass)
	if err != nil {
		return nil, err
	}
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Config.Port),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{keyPair},
			MinVersion:   tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP521,
				tls.CurveP384,
				tls.CurveP256,
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}
	return server, nil
}

// Start creates and starts the HTTPS server.
func Start() error {
	server, err := create()
	if err != nil {
		return err
	}
	log.Printf("server: listening on port %d\n", config.Config.Port)
	return server.ListenAndServeTLS("", "")
}
