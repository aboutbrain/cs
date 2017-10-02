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
	InputVectorSize            = 256
	OutputVectorSize           = 256
	ContextSize                = 6
	CombinatorialSpaceSize     = 60000
	ReceptorsPerPoint          = 32
	ClusterThreshold           = 6
	ClusterActivationThreshold = 4
	CharacterBits              = 8
	PointMemoryLimit           = 100
)

func main() {
	rand.Seed(time.Now().Unix())

	b, err := ioutil.ReadFile("TheOldManAndTheSea.txt") // just pass the file name
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

			fmt.Printf("i: %d, InputText  : \"%s\"\n", i*1000+j, textFragment)
			sourceCode := text.GetTextFragmentCode(textFragment, codes)

			targetText := txt5 + strings.Repeat("_", ContextSize - l)
			learningCode := text.GetTextFragmentCode(targetText, codes)
			fmt.Printf("i: %d, TargetText : \"%s\"\n", i*1000+j, targetText)
			//learningCode := wordCodeMap[word][0]

			mc.SetInputVector(sourceCode)
			mc.SetLearningVector(learningCode)

			mc.Next()
			if day == true {
				mc.AddNewClusters()
				fmt.Println("День")
			} else {
				fmt.Println("Ночь")
			}
			if t == 1000 {
				t = 0
				day = !day
			}

			total, permanent1, permanent2 := comSpace.ClustersCounters()
			outputVector := mc.OutVector()
			fmt.Printf("Clusters: %d, Permanent1: %d, Permanent2: %d\n", total, permanent1, permanent2)
			showVectors(sourceCode, outputVector, learningCode)

			comSpace.InternalTime++
			t++
		}
	}

	/*point := comSpace.Points[5]
	fmt.Printf("%#v\n", point)*/
}

func showVectors(source, output, learning bitarray.BitArray) {

	fmt.Printf("InputVector:   %s\n", cs.BitArrayToString(source, InputVectorSize))
	fmt.Printf("OutputVector:  %s\n", cs.BitArrayToString(output, OutputVectorSize))
	fmt.Printf("LerningVector: %s\n", cs.BitArrayToString(learning, OutputVectorSize))
	fmt.Printf("DeltaVector:   %s\n", BitArrayToString2(output, learning, OutputVectorSize))
	nVector := learning.Equals(output)
	if !nVector {
		fmt.Println("\033[31mFAIL!!\033[0m\n")
	} else {
		fmt.Println("\033[32mPASS!!\033[0m\n")
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
