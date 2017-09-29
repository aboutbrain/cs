package cs

import (
	"math/rand"
	"github.com/golang-collections/go-datastructures/bitarray"
)

func Random(min, max int) int {
	return rand.Intn(max - min) + min
}

func inArray(num int, arr []uint64) bool {
	for _, v := range arr {
		if int(v) == num {
			return true
		}
	}
	return false
}

func BitArrayToString(ba bitarray.BitArray) string {
	nums := ba.ToNums()
	s := ""
	for i := 0; i < 256; i++ {
		if inArray(i, nums) {
			s += "1"
		} else {
			s += "0"
		}
	}
	return s
}