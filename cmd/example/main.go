package main

import (
	"math/rand"
	"time"

	"fmt"

	"github.com/aboutbrain/cs"
	"github.com/aboutbrain/cs/persist"
	"github.com/aboutbrain/cs/text"
)

var _ = fmt.Printf // For debugging; delete when done.

const (
	InputVectorSize            = 128
	OutputVectorSize           = 128
	ContextSize                = 10
	CombinatorialSpaceSize     = 60000
	ReceptorsPerPoint          = 16
	OutputsPerPoint            = 16
	ClusterThreshold           = 6
	ClusterActivationThreshold = 4
	CharacterBits              = 8
	PointMemoryLimit           = 100
)

func main() {
	rand.Seed(time.Now().Unix())

	charContextVectors := text.GetCharContextMap(CharacterBits, text.Alpha, InputVectorSize, ContextSize)

	path := "codes.json"
	persist.ToFile(path, charContextVectors)
	codes := persist.FromFile(path)

	comSpace := cs.NewCombinatorialSpace(CombinatorialSpaceSize, ReceptorsPerPoint, OutputsPerPoint, OutputVectorSize)
	mc := cs.NewMiniColumn(ClusterThreshold, ClusterActivationThreshold, PointMemoryLimit, InputVectorSize, OutputVectorSize)
	mc.SetCombinatorialSpace(comSpace)

	day := true
	j := 0

	for i := 0; i < 10000; i += 1 {
		textFragment := text.GetTextFragment(sourceText, i, 1)
		//textFragment := "a"
		fmt.Printf("i: %d, InputText : \"%s\"\n", i, textFragment)
		sourceCode := text.GetTextFragmentCode(textFragment, codes)

		learningCode := codes.CharContext[int([]rune(textFragment)[0])][1]
		/*learningText := ""
		if i > 1 {
			learningText = text.GetTextFragment(i-1, 1)
		}
		fmt.Printf("i: %d, OutputText: \"%s\"\n", i, learningText)
		learningCode := text.GetTextFragmentCode(learningText, codes.CharContext)
		*/
		mc.SetInputVector(sourceCode)
		mc.SetLearningVector(learningCode)

		mc.Next()
		if day == true {
			mc.AddNewClusters()
			fmt.Println("День")
		} else {
			fmt.Println("Ночь")
		}
		if j == 300 {
			j = 0
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
		j++
	}

	point := comSpace.Points[5]
	fmt.Printf("%#v\n", point)
}
