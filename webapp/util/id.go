package util

import (
	"encoding/hex"

	"github.com/philippecery/libs/crng"
)

func GenerateUUID() string {
	rnd, err := crng.GetBytes(32)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(rnd)
}
