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
	InputVectorSize            = 256
	OutputVectorSize           = 256
	ContextSize                = 10
	CombinatorialSpaceSize     = 60000
	ReceptorsPerPoint          = 32
	OutputsPerPoint            = 10
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
	mc := cs.NewMiniColumn(ClusterThreshold, PointMemoryLimit)
	mc.SetCombinatorialSpace(comSpace)

	for i := 0; i < 100; i += 5 {
		textFragment := text.GetTextFragment(i, 10)
		fmt.Printf("TextFragment: \"%s\"\n", textFragment)
		sourceCode := text.GetTextFragmentCode(textFragment, codes.CharContext)

		learningCode := sourceCode

		mc.SetInputVector(sourceCode)
		mc.SetLearningVector(learningCode)

		mc.Next()
		mc.AddNewClusters()

		fmt.Printf("Clusters: %d\n", comSpace.GetClustersCounter())
	}

	point := comSpace.Points[5]
	fmt.Printf("%#v\n", point)
}
