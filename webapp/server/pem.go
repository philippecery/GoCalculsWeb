package server

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/philippecery/maths/webapp/util"

	"github.com/youmark/pkcs8"
)

// KeyPair contains the server's certificate and private key.
type KeyPair struct {
	PrivateKey  []byte
	Certificate []byte
}

func loadKeyPair(pemFilePath string, keyPassword *[]byte) (*KeyPair, error) {
	defer util.Clear(keyPassword)
	remainder, err := ioutil.ReadFile(pemFilePath)
	if err != nil {
		return nil, err
	}
	var keyPair = &KeyPair{}
	var pemBlock *pem.Block
	for len(remainder) > 0 {
		pemBlock, remainder = pem.Decode(remainder)
		if pemBlock != nil {
			if pemBlock.Type == "CERTIFICATE" {
				if keyPair.Certificate == nil {
					keyPair.Certificate = loadCertificate(pemBlock)
				} else {
					return nil, fmt.Errorf("ssl: Multiple certificates found in %s", pemFilePath)
				}
			} else {
				if keyPassword != nil {
					if keyPair.PrivateKey == nil {
						var pk interface{}
						if pk, err = pkcs8.ParsePKCS8PrivateKey(pemBlock.Bytes, *keyPassword); err != nil {
							return nil, err
						}
						var keyDER []byte
						if keyDER, err = pkcs8.MarshalPrivateKey(pk, nil, nil); err != nil {
							return nil, err
						}
						keyPair.PrivateKey = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
					} else {
						return nil, fmt.Errorf("ssl: Multiple private keys found in %s", pemFilePath)
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
