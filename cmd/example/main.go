package main

import (
	"math/rand"
	"time"

	"fmt"

	"github.com/aboutbrain/cs"
	"github.com/aboutbrain/cs/text"
)

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

	codes := text.GetCharContextMap(CharacterBits, text.Alpha, InputVectorSize, ContextSize)
	//fmt.Printf("%+v\n", codes)

	//path := "aaa.json"

	//persist.ToFile(path, codes)

	/*codes2 := persist.FromFile(path)
	fmt.Printf("%+v\n", codes2)*/

	textFragment := text.GetTextFragment(0, 5)
	sourceCode := text.GetTextFragmentCode(textFragment, codes.CharContext)
	//fmt.Printf("%+v\n", sourceCode)

	comSpace := cs.NewCombinatorialSpace(CombinatorialSpaceSize, ReceptorsPerPoint, OutputsPerPoint, OutputVectorSize)

	//fmt.Printf("%+v\n", textFragment)

	learningCode := sourceCode

	mc := cs.NewMiniColumn(ClusterThreshold, PointMemoryLimit)
	mc.SetCombinatorialSpace(comSpace)

	mc.SetInputVector(sourceCode)
	mc.SetLearningVector(learningCode)
	mc.Next()

	mc.AddNewClusters()

	fmt.Println(comSpace.GetClustersCounter())

	point := comSpace.Points[5]
	fmt.Printf("%+v\n", point)
}
