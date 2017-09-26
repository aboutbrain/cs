package cs

import (
	"math/rand"
)

func Random(min, max int) int {
	return rand.Intn(max - min) + min
}
