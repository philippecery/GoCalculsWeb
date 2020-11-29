package main

import (
	"crypto/rand"
	hash "crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	password := os.Args[1]
	hashedPwd := make([]byte, 0)
	salt := make([]byte, 32)
	rand.Read(salt)
	h := hash.New()
	h.Write(salt)
	h.Write([]byte(password))
	hashedPwd = append(hashedPwd, salt...)
	hashedPwd = append(hashedPwd, h.Sum(nil)...)
	fmt.Println(base64.StdEncoding.EncodeToString(hashedPwd))
}
