package cs

import (
	"math/rand"
)

func random(min, max int) int {
	return rand.Intn(max - min) + min
}
