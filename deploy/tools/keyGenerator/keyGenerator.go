package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"

	"github.com/philippecery/libs/crng"
)

func main() {
	if len(os.Args) == 2 {
		arg1 := os.Args[1]
		if length, err := strconv.Atoi(arg1); err == nil {
			fmt.Printf("Key: %s\n", base64.StdEncoding.EncodeToString(crng.GetBytes(length)))
		} else {
			fmt.Printf("Invalid length: %s\n", arg1)
		}
	}
}
