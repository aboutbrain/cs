package main

import (
	"fmt"

	"github.com/aboutbrain/cs"

	"math/rand"
	"time"

	"encoding/json"

	"github.com/aboutbrain/cs/persist"
	"github.com/aboutbrain/cs/text"
)

const (
	InputVectorSize            = 256
	OutputVectorSize           = 256
	ContextSize                = 10
	CombinatorialSpaceSize     = 60000
	BitsPerPoint               = 32
	ClusterThreshold           = 6
	ClusterActivationThreshold = 4
	CharacterBits              = 8
	PointMemoryLimit           = 100
)

func main() {
	rand.Seed(time.Now().Unix())

	const alpha = " abcdefghijklmnopqrstuvwxyz"
	charCodes := text.GetCharCodes(CharacterBits, alpha, InputVectorSize, 0, 127)
	fmt.Printf("%+v", charCodes)
	contextCodes := text.GetContextCodes(ContextSize, InputVectorSize, 128, 255)
	fmt.Printf("%+v", contextCodes)

	charContextCodes := text.GetCharContextCodes(charCodes, contextCodes)
	fmt.Printf("%+v", charContextCodes)

	capCode := charContextCodes[32][0].Capacity()
	fmt.Println(capCode)

	j, _ := json.Marshal(charCodes)
	fmt.Printf("%+v", j)

	path := "aaa.json"

	persist.ToFile(path, &charContextCodes)

	codes2 := persist.FromFile(path)

	textFragment := text.GetTextFragment(0, 1)
	sourceCode := text.GetTextFragmentCode(textFragment, *codes2)
	fmt.Printf("%+v", sourceCode)

	comSpace := cs.NewCombinatorialSpace(CombinatorialSpaceSize, BitsPerPoint, OutputVectorSize)

	learningCode := sourceCode

	mc := cs.NewMiniColumn(ClusterThreshold, PointMemoryLimit)
	mc.SetCombinatorialSpace(comSpace)

	mc.SetInputVector(sourceCode)
	mc.SetLearningVector(learningCode)
	mc.Next()

	/*mc.ActivateClusters()
	mc.MakeOutVector()
	mc.ModifyClusters()
	mc.ConsolidateMemory()*/

	mc.AddNewClusters()

	fmt.Println(mc)
}
