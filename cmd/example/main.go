package main

import (
	"fmt"

	"strings"
	"unicode"

	"io/ioutil"

	"math/rand"
	"time"

	"github.com/aboutbrain/cs"
	"github.com/aboutbrain/cs/persist"
	"github.com/aboutbrain/cs/text"
	"github.com/golang-collections/go-datastructures/bitarray"
)

var _ = fmt.Printf // For debugging; delete when done.

const (
	InputVectorSize            = 128
	OutputVectorSize           = 128
	ContextSize                = 6
	CombinatorialSpaceSize     = 5000
	ReceptorsPerPoint          = 32
	ClusterThreshold           = 5
	ClusterActivationThreshold = 4
	CharacterBits              = 8
	PointMemoryLimit           = 100
)

func main() {
	rand.Seed(time.Now().Unix())

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

	charContextVectors := text.GetCharContextMap(CharacterBits, text.Alpha, InputVectorSize, ContextSize)

	path := "codes.json"
	persist.ToFile(path, charContextVectors)
	codes := persist.FromFile(path)

	comSpace := cs.NewCombinatorialSpace(CombinatorialSpaceSize, InputVectorSize, ReceptorsPerPoint, OutputVectorSize)
	mc := cs.NewMiniColumn(ClusterThreshold, ClusterActivationThreshold, PointMemoryLimit, InputVectorSize, OutputVectorSize)
	mc.SetCombinatorialSpace(comSpace)

	day := true
	t := 0

	//fragmentLength := 5
	for i := 0; i < 100; i++ {
		for j := 0; j < 1000; j++ {
			word := strings.ToLower(words[j])
			l := len(word)
			if l >= 5 {
				l = 5
			}

			txt5 := word[:l]
			wordContextSize := ContextSize - l
			context := cs.Random(0, wordContextSize)

			//context := 0
			textFragment := strings.Repeat("_", context)
			textFragment += txt5
			after := strings.Repeat("_", ContextSize-len(textFragment))
			textFragment += after

			sourceCode := text.GetTextFragmentCode(textFragment, codes)
			inputBits := len(sourceCode.ToNums())
			fmt.Printf("i: %d, InputText  : \"%s\", Bit: %d\n", i*1000+j, textFragment, inputBits)

			targetText := txt5 + strings.Repeat("_", ContextSize-l)
			learningCode := text.GetTextFragmentCode(targetText, codes)
			lerningBits := len(learningCode.ToNums())
			fmt.Printf("i: %d, TargetText : \"%s\", Bit: %d\n", i*1000+j, targetText, lerningBits)
			//learningCode := wordCodeMap[word][0]

			mc.SetInputVector(sourceCode)
			mc.SetLearningVector(learningCode)

			outputVector := mc.Calculate()

			nVector := learningCode.Equals(outputVector)

			if day == true && !nVector {
				s := "Day"
				if !nVector {
					s += " - learning!"
					mc.Learn(day)
				}
				fmt.Println(s)
			} else {
				s := "Night"
				if nVector {
					s += " - learned!"
				}
				fmt.Println(s)
			}
			if t == 1000 {
				t = 0
				day = !day
			}

			total, permanent1, permanent2 := comSpace.ClustersCounters()

			showVectors(sourceCode, outputVector, learningCode, nVector)
			fmt.Printf("Clusters: %d, Permanent1: %d, Permanent2: %d\n\n", total, permanent1, permanent2)

			comSpace.InternalTime++
			t++
		}
	}
}

func showVectors(source, output, learning bitarray.BitArray, nVector bool) {

	fmt.Printf("InputVector:   %s\n", cs.BitArrayToString(source, InputVectorSize))
	fmt.Printf("OutputVector:  %s\n", cs.BitArrayToString(output, OutputVectorSize))
	fmt.Printf("LerningVector: %s\n", cs.BitArrayToString(learning, OutputVectorSize))
	fmt.Printf("DeltaVector:   %s\n", BitArrayToString2(output, learning, OutputVectorSize))

	if !nVector {
		fmt.Println("\033[31mFAIL!!\033[0m")
	} else {
		fmt.Println("\033[32mPASS!!\033[0m")
	}
}

func BitArrayToString2(output, learning bitarray.BitArray, vectorLen int) string {
	delta := output.And(learning)
	nums := delta.ToNums()
	s := ""
	for i := 0; i < vectorLen; i++ {
		if cs.InArray(i, nums) {
			s += "\033[32m1\033[0m"
		} else {
			s += "0"
		}
	}
	return s
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
