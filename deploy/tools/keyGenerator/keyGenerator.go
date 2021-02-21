package main

import (
	"encoding/base64"
	"fmt"

	"github.com/philippecery/libs/crng"
)

func main() {
	fmt.Printf("Key: %s\n", base64.StdEncoding.EncodeToString(crng.GetBytes(32)))
}
