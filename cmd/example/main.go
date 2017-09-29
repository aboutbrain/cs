package main

import (
	"math/rand"
	"time"

	"fmt"

	"github.com/aboutbrain/cs"
	"github.com/aboutbrain/cs/persist"
	"github.com/aboutbrain/cs/text"
	"github.com/golang-collections/go-datastructures/bitarray"
)

var _ = fmt.Printf // For debugging; delete when done.

const (
	InputVectorSize            = 256
	OutputVectorSize           = 256
	ContextSize                = 10
	CombinatorialSpaceSize     = 60000
	ReceptorsPerPoint          = 32
	OutputsPerPoint            = 32
	ClusterThreshold           = 6
	ClusterActivationThreshold = 4
	CharacterBits              = 8
	PointMemoryLimit           = 100
)

func main() {
	rand.Seed(time.Now().Unix())

	//charContextVectors := text.GetCharContextMap(CharacterBits, text.Alpha, InputVectorSize, ContextSize)

	path := "codes.json"
	//persist.ToFile(path, charContextVectors)
	codes := persist.FromFile(path)

	comSpace := cs.NewCombinatorialSpace(CombinatorialSpaceSize, ReceptorsPerPoint, OutputsPerPoint, OutputVectorSize)
	mc := cs.NewMiniColumn(ClusterThreshold, PointMemoryLimit)
	mc.SetCombinatorialSpace(comSpace)

	for i := 0; i < 100; i += 1 {
		textFragment := text.GetTextFragment(i, 1)
		//textFragment := "a"
		fmt.Printf("TextFragment: \"%s\"\n", textFragment)
		sourceCode := text.GetTextFragmentCode(textFragment, codes.CharContext)

		learningCode := codes.CharContext[int([]rune(textFragment)[0])][1]

		mc.SetInputVector(sourceCode)
		mc.SetLearningVector(learningCode)

		mc.Next()
		mc.AddNewClusters()

		fmt.Printf("Clusters: %d\n", comSpace.GetClustersCounter())
		fmt.Printf("InputVector:   %s\n", BitArrayToString(sourceCode))
		fmt.Printf("OutputVector:  %s\n", BitArrayToString(mc.OutVector()))
		fmt.Printf("LerningVector: %s\n", BitArrayToString(learningCode))
		nVector := learningCode.Equals(mc.OutVector())
		if !nVector {
			fmt.Println("\033[31mFAIL!!\033[0m")
		} else {
			fmt.Println("\033[32mPASS!!\033[0m")
		}
	}

	point := comSpace.Points[5]
	fmt.Printf("%#v\n", point)
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

func inArray(num int, arr []uint64) bool {
	for _, v := range arr {
		if int(v) == num {
			return true
		}
	}
	return false
}
