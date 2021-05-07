package util

import (
	"crypto/tls"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/philippecery/libs/bytes"

	"github.com/youmark/pkcs8"
)

type keyPair struct {
	privateKey  []byte
	certificate []byte
}

func X509KeyPair(keystorePath string, serverKeyPass *[]byte) (tls.Certificate, error) {
	var nilCert tls.Certificate
	kp, err := loadKeyPair(keystorePath, serverKeyPass)
	if err != nil {
		return nilCert, err
	}
	if kp.certificate == nil {
		return nilCert, errors.New("util: no certificate found")
	}
	if kp.privateKey == nil {
		return nilCert, errors.New("util: no private key found")
	}
	return tls.X509KeyPair(kp.certificate, kp.privateKey)
}

func loadKeyPair(pemFilePath string, keyPassword *[]byte) (*keyPair, error) {
	defer bytes.Clear(keyPassword)
	remainder, err := ioutil.ReadFile(pemFilePath)
	if err != nil {
		return nil, err
	}
	var keyPair = &keyPair{}
	var pemBlock *pem.Block
	for len(remainder) > 0 {
		pemBlock, remainder = pem.Decode(remainder)
		if pemBlock != nil {
			if pemBlock.Type == "CERTIFICATE" {
				if keyPair.certificate == nil {
					keyPair.certificate = loadCertificate(pemBlock)
				} else {
					return nil, fmt.Errorf("util: Multiple certificates found in %s", pemFilePath)
				}
			} else {
				if keyPassword != nil {
					if keyPair.privateKey == nil {
						var pk interface{}
						if pk, err = pkcs8.ParsePKCS8PrivateKey(pemBlock.Bytes, *keyPassword); err != nil {
							return nil, err
						}
						var keyDER []byte
						if keyDER, err = pkcs8.MarshalPrivateKey(pk, nil, nil); err != nil {
							return nil, err
						}
						keyPair.privateKey = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
					} else {
						return nil, fmt.Errorf("util: Multiple private keys found in %s", pemFilePath)
					}
				}
			}
		}
	}
	return keyPair, nil
}

func loadCertificate(certBlock *pem.Block) []byte {
	return pem.EncodeToMemory(certBlock)
}
