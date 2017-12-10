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

func BitArrayToString(ba bitarray.BitArray, vectorLen int) (string, int) {
	nums := ba.ToNums()
	s := ""
	for i := 0; i < vectorLen; i++ {
		if InArray64(vectorLen-1-i, nums) {
			s += "1"
		} else {
			s += "0"
		}
	}
	return s, len(nums)
}

func RotateL(s string, i int) string {
	a := []byte(s)
	x, b := (a)[:i], (a)[i:]
	a = append(b, x...)
	return string(a)
}

func RotateR(a *[]byte, i int) {
	x, b := (*a)[:(len(*a)-i)], (*a)[(len(*a)-i):]
	*a = append(b, x...)
}

func MinMax(x map[int]votingCounters) (int, int) {
	var n, smallest, biggest int
	/*x := []int{
		48,96,86,68,
		57,82,63,70,
		37,34,83,27,
		19,97, 9,17,
	}*/

	for _, v := range x {
		if v.Up > n {
			//fmt.Println(v, ">", n)
			n = v.Up
			biggest = n
		} else {
			//fmt.Println(v, "<", n)
		}
	}

	//fmt.Println("The biggest number is ", biggest)
	for _, v := range x {
		if v.Up > n {
			//fmt.Println(v, ">", n)
		} else {
			//fmt.Println(v, "<", n)
			n = v.Up
			smallest = n
		}
	}
	//fmt.Println("The smallest number is ", smallest)
	return smallest, biggest
}
