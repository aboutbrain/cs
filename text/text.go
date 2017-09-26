package text

import (
	"strings"

	"github.com/aboutbrain/cs"
	"github.com/golang-collections/go-datastructures/bitarray"
)

const Alpha = " abcdefghijklmnopqrstuvwxyz"

//const Alpha = " ABCDEFGHIJKLMNOPQRSTUVWXYZ"

//isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString

type CharContext map[int]map[int]bitarray.BitArray

type CharContextCodes struct {
	VectorCapacity int
	CharContext    CharContext
}

func GetCharContextMap(bitPerChar int, alpha string, capacity int, contextSize int) *CharContextCodes {
	codes := make(CharContext)
	for _, char := range alpha {
		m := make(map[int]bitarray.BitArray)
		codes[int(char)] = m
		for i := 0; i < contextSize; i++ {
			arr := bitarray.NewBitArray(uint64(capacity))
			for j := 0; j < bitPerChar; j++ {
			rnd:
				bitNumber := cs.Random(0, capacity-1)
				if a, _ := arr.GetBit(uint64(bitNumber)); a != true {
					arr.SetBit(uint64(bitNumber))
				} else {
					goto rnd
				}
			}
			codes[int(char)][i] = arr
		}
	}
	charContextCode := CharContextCodes{VectorCapacity: capacity, CharContext: codes}
	return &charContextCode
}

func GetTextFragment(start, end int) string {
	text := `I proffered him my passport and stood the suitcase on the
white counter. The inspector rapidly leafed through it with his
long careful fingers. He was dressed in a  white  uniform  with
silver  buttons  and silver braid on the shoulders. He laid the
passport aside and touched the suitcase with the  tips  of  his
fingers.`
	return strings.ToLower(text[start:end])
}

func GetTextFragmentCode(txtFragment string, charContextCodes CharContext) bitarray.BitArray {
	code := bitarray.NewBitArray(256)
	for i, char := range txtFragment {
		codeCurrent := charContextCodes[int(char)][i]
		code = code.Or(codeCurrent)
	}
	return code
}
