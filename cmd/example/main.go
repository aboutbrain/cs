package main

import (
	"fmt"

	"github.com/aboutbrain/cs"
	"github.com/golang-collections/go-datastructures/bitarray"

	"math/rand"
	"time"
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
	comSpace := cs.NewCombinatorialSpace(CombinatorialSpaceSize, BitsPerPoint, OutputVectorSize)

	sourceCode := bitarray.NewBitArray(InputVectorSize)
	sourceCode.SetBit(0)
	sourceCode.SetBit(2)
	sourceCode.SetBit(3)
	sourceCode.SetBit(5)
	sourceCode.SetBit(15)
	sourceCode.SetBit(21)
	sourceCode.SetBit(25)
	sourceCode.SetBit(35)
	sourceCode.SetBit(45)
	sourceCode.SetBit(63)

	learningCode := bitarray.NewBitArray(OutputVectorSize)
	learningCode.SetBit(0)
	learningCode.SetBit(2)
	learningCode.SetBit(3)
	learningCode.SetBit(63)

	mc := cs.NewMiniColumn(ClusterThreshold, PointMemoryLimit)
	mc.SetCombinatorialSpace(comSpace)

	mc.SetInputVector(sourceCode)
	mc.SetLearningVector(learningCode)
	mc.Next()

	/*mc.ActivateClusters()
	mc.MakeOutVector()
	mc.ModyfyClusters()
	mc.ConsolidateMemory()*/

	mc.AddNewClusters()

	fmt.Println(mc)
}
