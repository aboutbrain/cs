package main

import (
	"fmt"

	"strings"
	"unicode"

	"io/ioutil"

	"github.com/aboutbrain/cs"
	"github.com/aboutbrain/cs/persist"
	"github.com/aboutbrain/cs/text"
	"github.com/golang-collections/go-datastructures/bitarray"
)

var _ = fmt.Printf // For debugging; delete when done.

const (
	InputVectorSize            = 256
	OutputVectorSize           = 256
	ContextSize                = 15
	CombinatorialSpaceSize     = 60000
	ReceptorsPerPoint          = 32
	ClusterThreshold           = 6
	ClusterActivationThreshold = 4
	CharacterBits              = 8
	PointMemoryLimit           = 100
)

func main() {
	//rand.Seed(time.Now().Unix())

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

	wordCodeMap := make(map[string]bitarray.BitArray)

	for _, word := range words {
		wordCodeMap[strings.ToLower(word)] = getRandomCode(16, OutputVectorSize)
	}
	fmt.Printf("WordsCount: %d\n", len(wordCodeMap))

	charContextVectors := text.GetCharContextMap(CharacterBits, text.Alpha, InputVectorSize, ContextSize)

	path := "codes.json"
	persist.ToFile(path, charContextVectors)
	codes := persist.FromFile(path)

	comSpace := cs.NewCombinatorialSpace(CombinatorialSpaceSize, ReceptorsPerPoint, OutputVectorSize)
	mc := cs.NewMiniColumn(ClusterThreshold, ClusterActivationThreshold, PointMemoryLimit, InputVectorSize, OutputVectorSize)
	mc.SetCombinatorialSpace(comSpace)

	day := true
	t := 0

	//fragmentLength := 5
	for i := 0; i < 10; i++ {
		for j := 0; j < 1000; j++ {
			textFragment := strings.ToLower(words[j])
			fmt.Printf("i: %d, InputText : \"%s\"\n", j, textFragment)
			sourceCode := text.GetTextFragmentCode(textFragment, codes)

			learningCode := wordCodeMap[textFragment]

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

			total, permanent := comSpace.ClustersCounters()
			fmt.Printf("Clusters: %d, Permanent: %d\n", total, permanent)
			fmt.Printf("InputVector:   %s\n", cs.BitArrayToString(sourceCode, InputVectorSize))
			fmt.Printf("OutputVector:  %s\n", cs.BitArrayToString(mc.OutVector(), OutputVectorSize))
			fmt.Printf("LerningVector: %s\n", cs.BitArrayToString(learningCode, OutputVectorSize))
			nVector := learningCode.Equals(mc.OutVector())
			if !nVector {
				fmt.Println("\033[31mFAIL!!\033[0m\n")
			} else {
				fmt.Println("\033[32mPASS!!\033[0m\n")
			}
			comSpace.InternalTime++
			t++
		}
	}

	point := comSpace.Points[5]
	fmt.Printf("%#v\n", point)
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
