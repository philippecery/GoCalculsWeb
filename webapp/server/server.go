package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/config"
)

func create() (*http.Server, error) {
	if config.Config.SSL == nil {
		return nil, errors.New("Missing SSL configuration")
	}
	if config.Config.SSL.Password == "" {
		return nil, errors.New("Missing private key password")
	}
	var serverKeyPass = []byte(config.Config.SSL.Password)
	kp, err := loadKeyPair(config.Config.SSL.Keystore, &serverKeyPass)
	if err != nil {
		return nil, err
	}
	if kp.Certificate == nil {
		return nil, errors.New("No certificate found")
	}
	if kp.PrivateKey == nil {
		return nil, errors.New("No private key found")
	}
	keyPair, err := tls.X509KeyPair(kp.Certificate, kp.PrivateKey)
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
	log.Printf("Starting server, listening on port %d\n", config.Config.Port)
	return server.ListenAndServeTLS("", "")
}
