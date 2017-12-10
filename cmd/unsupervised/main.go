package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/aboutbrain/cs"
	"github.com/aboutbrain/cs/bitarray"
	"github.com/aboutbrain/cs/persist"
)

var _ = fmt.Printf // For debugging; delete when done.

const (
	InputVectorSize            = 128
	OutputVectorSize           = 32
	ContextSize                = 5
	CombinatorialSpaceSize     = 60000
	ReceptorsPerPoint          = 24
	OutputsPerPoint            = 8
	ClusterThreshold           = 6
	ClusterActivationThreshold = 4
	CharacterBits              = 5
	PointMemoryLimit           = 100
	Level                      = 10 //5
)

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()
	//rand.Seed(time.Now().Unix())

	b, err := ioutil.ReadFile("./testdata/TheOldManAndTheSea.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	str := string(b) // convert content to a 'string'

	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	words := strings.FieldsFunc(str, f)
	fmt.Printf("WordsTotal: %d\n", len(words))

	wordCodeMap := make(map[string]map[int]bitarray.BitArray)

	for _, word := range words {
		l := len(word)
		wordContextSize := ContextSize - l
		m := make(map[int]bitarray.BitArray)
		for i := 0; i < wordContextSize; i++ {
			m[i] = getRandomCode(16, OutputVectorSize)
		}
		wordCodeMap[strings.ToLower(word)] = m
	}
	fmt.Printf("WordsCount: %d\n", len(wordCodeMap))

	charContextVectors := cs.GetCharContextMap(CharacterBits, cs.Alpha, InputVectorSize, ContextSize)

	path := "codes.json"
	persist.ToFile(path, charContextVectors)
	codes := persist.FromFile(path)

	comSpace := cs.NewCombinatorialSpace(CombinatorialSpaceSize, InputVectorSize, ReceptorsPerPoint, OutputsPerPoint, OutputVectorSize)
	mc := cs.NewMiniColumn(ClusterThreshold, ClusterActivationThreshold, PointMemoryLimit, InputVectorSize, OutputVectorSize, Level)
	mc.SetCombinatorialSpace(comSpace)

	s := ""
	for _, word := range words {
		s += word + "_"
	}

	//textPosition := 0
	offset := 0
	//Fragment := 14

	cortex := cs.NewCortex(mc, 5, codes)

	for wordId, word := range words {
		//c := 0
		if len(word) > 1 {
			word = word[0:1]
		}
		//word = "abcfd"
		txt := strings.ToLower(word)
		textFragment := strings.Repeat("_", offset)
		textFragment += txt

		i2 := ContextSize - len(textFragment)
		if i2 > 0 {
			after := strings.Repeat("_", i2)
			textFragment += after
		}

		//textFragment = cs.RotateL(textFragment, cs.Random(0, 5))

		contextData := make([]cs.CortexContext, ContextSize)

		for i := 0; i < ContextSize; i++ {
			textFragmentContext := cs.RotateL(textFragment, i)
			sourceCode := cs.GetTextFragmentCode(textFragmentContext, codes)
			contextData[i] = cs.CortexContext{
				TextFragment: textFragmentContext,
				InputVector:  sourceCode,
				Id:           wordId,
			}
		}

		cortex.Run(&contextData)

	}
}

func getRandomCode(bitPerWord, capacity int) bitarray.BitArray {
	arr := bitarray.NewBitArray(uint64(capacity))
	for j := 0; j < bitPerWord; j++ {
	rnd:
		bitNumber := cs.Random(0, capacity-1)
		if a, _ := arr.GetBit(uint64(bitNumber)); a != true {
			arr.SetBit(uint64(bitNumber))
		} else {
			goto rnd
		}
	}
	return arr
}
