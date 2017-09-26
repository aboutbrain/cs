package text

import (
	"fmt"
	"strings"

	"github.com/aboutbrain/cs"
	"github.com/golang-collections/go-datastructures/bitarray"
)

const alpha = " abcdefghijklmnopqrstuvwxyz"

//const Alpha = " ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetCharCodes(bitPerChar int, alpha string, capacity int, startBit, stopBit int) map[int]bitarray.BitArray {
	codes := make(map[int]bitarray.BitArray)
	for _, char := range alpha {
		charVector := bitarray.NewBitArray(uint64(capacity))
		for j := 0; j < bitPerChar; j++ {
			bitNumber := cs.Random(startBit, stopBit)
			charVector.SetBit(uint64(bitNumber))
		}
		codes[int(char)] = charVector
	}
	return codes
}

func GetContextCodes(contextSize int, capacity int, startBit, stopBit int) []bitarray.BitArray {
	codes := make([]bitarray.BitArray, contextSize)
	for i := 0; i < contextSize; i++ {
		contextVector := bitarray.NewBitArray(uint64(capacity))
		for j := 0; j < 8; j++ {
			bitNumber := cs.Random(startBit, stopBit)
			contextVector.SetBit(uint64(bitNumber))
		}
		codes[i] = contextVector
	}
	return codes
}

type ContextArray map[int]bitarray.BitArray
type CharContextMap map[int]ContextArray

func GetCharContextCodes(charCodes map[int]bitarray.BitArray, contextCodes []bitarray.BitArray) CharContextMap {
	codes := make(CharContextMap)
	for i, charCode := range charCodes {
		if charCode != nil {
			contArray := ContextArray{}
			for j, contextCode := range contextCodes {
				resultCode := charCode.Or(contextCode)
				contArray[j] = resultCode
				codes[int(i)] = contArray
				fmt.Print(i, j, contextCode, resultCode)
			}
		}
	}
	return codes
}

func GetTextFragment(start, end int) string {
	text := `I  proffered him my passport and stood the suitcase on the
white counter. The inspector rapidly leafed through it with his
long careful fingers. He was dressed in a  white  uniform  with
silver  buttons  and silver braid on the shoulders. He laid the
passport aside and touched the suitcase with the  tips  of  his
fingers.`
	return strings.ToLower(text[start:end])
}

func GetTextFragmentCode(txtFragment string, charContextCodes CharContextMap) bitarray.BitArray {
	code := bitarray.NewBitArray(256)
	for i, char := range txtFragment {
		codeCurrent := charContextCodes[int(char)][i]
		code = code.Or(codeCurrent)
	}
	return code
}
