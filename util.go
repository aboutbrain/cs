package cs

import (
	"math/rand"

	"github.com/aboutbrain/cs/bitarray"
)

func Random(min, max int) int {
	return rand.Intn(max-min) + min
}

func InArray8(num int, arr []uint8) bool {
	for _, v := range arr {
		if int(v) == num {
			return true
		}
	}
	return false
}

func InArray64(num int, arr []uint64) bool {
	for _, v := range arr {
		if int(v) == num {
			return true
		}
	}
	return false
}

func BitArrayToString(ba bitarray.BitArray, vectorLen int) string {
	nums := ba.ToNums()
	s := ""
	for i := 0; i < vectorLen; i++ {
		if InArray64(vectorLen-1-i, nums) {
			s += "1"
		} else {
			s += "0"
		}
	}
	return s
}
