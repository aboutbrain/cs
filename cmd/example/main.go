package main

import (
	"fmt"

	"io/ioutil"
	"strings"
	"unicode"

	"github.com/aboutbrain/cs"
	"github.com/aboutbrain/cs/bitarray"
	"github.com/aboutbrain/cs/persist"
	"github.com/aboutbrain/cs/text"
)

var _ = fmt.Printf // For debugging; delete when done.

const (
	InputVectorSize            = 128
	OutputVectorSize           = 128
	ContextSize                = 10
	CombinatorialSpaceSize     = 20000
	ReceptorsPerPoint          = 24
	OutputsPerPoint            = 24
	ClusterThreshold           = 6
	ClusterActivationThreshold = 4
	CharacterBits              = 5
	PointMemoryLimit           = 100
	Level                      = 100
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

	charContextVectors := text.GetCharContextMap(CharacterBits, text.Alpha, InputVectorSize, ContextSize)

	path := "codes.json"
	persist.ToFile(path, charContextVectors)
	codes := persist.FromFile(path)

	comSpace := cs.NewCombinatorialSpace(CombinatorialSpaceSize, InputVectorSize, ReceptorsPerPoint, OutputsPerPoint, OutputVectorSize)
	mc := cs.NewMiniColumn(ClusterThreshold, ClusterActivationThreshold, PointMemoryLimit, InputVectorSize, OutputVectorSize, Level)
	mc.SetCombinatorialSpace(comSpace)

	day := true
	t := 0

	s := ""
	for _, word := range words {
		s += word + "_"
	}

	memOn := false

	textPosition := 0
	offset := 0

	const Segment = 3000
	const DayTime = 500
	Fragment := 9

	for i := 0; i < 10; i++ {
		for j := 0; j < Segment; j++ {
			context := 1

			/*if Fragment < 2 {
				offset = cs.Random(0, 6-Fragment)
			} else {
				offset = 0
			}*/

			txt := strings.ToLower(s[textPosition : textPosition+Fragment])

			textFragment := strings.Repeat("_", offset)
			textFragment += txt
			after := strings.Repeat("_", ContextSize-len(textFragment))
			textFragment += after

			sourceCode := text.GetTextFragmentCode(textFragment, codes)
			inputBits := mc.SetInputVector(sourceCode)
			mc.InputText = textFragment
			//inputBits := len(sourceCode.ToNums())
			fmt.Printf("i: %d, InputText  : \"%s\", Bit: %d\n", i*Segment+j, textFragment, inputBits)

			targetText := strings.Repeat("_", offset+context) + txt
			targetText += strings.Repeat("_", ContextSize-len(targetText))
			learningCode := text.GetTextFragmentCode(targetText, codes)
			learningBits := mc.SetLearningVector(learningCode)
			//learningBits := len(learningCode.ToNums())
			fmt.Printf("i: %d, TargetText : \"%s\", Bit: %d\n", i*Segment+j, targetText, learningBits)
			//learningCode := wordCodeMap[word][0]

			outputVector := mc.Calculate()

			nVector := learningCode.Equals(outputVector)

			memOn = false
			if day == true {
				s := "Day"
				if !nVector {
					s += " - learning!"
					memOn = true
				}
				fmt.Println(s)
			} else {
				//memOn = false
				s := "Night"
				if nVector {
					s += " - learned!"
				}
				fmt.Println(s)
			}
			mc.Learn(memOn)
			if t == DayTime {
				t = 0
				day = !day
			}

			total, permanent1, permanent2 := comSpace.ClustersCounters()

			showVectors(sourceCode, outputVector, learningCode, nVector)
			fmt.Printf("Clusters: %d, Permanent-1: %d, Permanent-2: %d\n\n", total, permanent1, permanent2)

			comSpace.InternalTime++
			t++
			textPosition++
		}
		/*if !day {
			Fragment++
		}*/
	}
}

func showVectors(source, output, learning bitarray.BitArray, nVector bool) {

	fmt.Printf("InputVector:   %s\n", cs.BitArrayToString(source, InputVectorSize))
	fmt.Printf("OutputVector:  %s\n", cs.BitArrayToString(output, OutputVectorSize))
	fmt.Printf("LerningVector: %s\n", cs.BitArrayToString(learning, OutputVectorSize))
	fmt.Printf("DeltaVector:   %s\n", BitArrayToString2(output, learning, OutputVectorSize))

	if len(output.ToNums()) > 0 {
		if !nVector {
			fmt.Println("\033[31mFAIL!!\033[0m")
		} else {
			fmt.Println("\033[32mPASS!!\033[0m")
		}
	}
}

func BitArrayToString2(output, learning bitarray.BitArray, vectorLen int) string {
	delta := output.And(learning)
	nums := delta.ToNums()
	s := ""
	for i := 0; i < vectorLen; i++ {
		if cs.InArray64(vectorLen-1-i, nums) {
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
