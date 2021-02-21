package main

import (
	"fmt"
	"log"
	"os"

	"github.com/philippecery/maths/webapp/util"
)

func main() {
	if len(os.Args) == 4 {
		what := os.Args[2]
		data := os.Args[3]
		var result string
		switch what {
		case "-password":
			result = util.ProtectPassword(data)
		case "-userid":
			result, _ = util.ProtectUserID(data)
		case "-pii":
			result, _ = util.ProtectPII(data)
		}
		fmt.Println(result)
		return
	}
	log.Fatalf("Invalid number of parameters: %d", len(os.Args))
}
