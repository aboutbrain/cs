package main

import (
	"fmt"

	"github.com/aboutbrain/cs"
	"github.com/golang-collections/go-datastructures/bitarray"

	"math/rand"
	"time"
)

const (
	InputVectorSize            = 8
	OutputVectorSize           = 8
	ContextSize                = 10
	CombinatorialSpaceSize     = 10
	BitsPerPoint               = 8
	ClusterThreshold           = 6
	ClusterActivationThreshold = 4
	CharacterBits              = 8
	PointMemoryCapacity        = 10
	PointContextCapacity       = 10
)

func main() {
	rand.Seed(time.Now().Unix())
	comSpace := cs.NewCombinatorialSpace(CombinatorialSpaceSize, BitsPerPoint)
	point := comSpace.Points[0]

	sourceCode := bitarray.NewBitArray(InputVectorSize)
	//fmt.Println(sourceCode)
	sourceCode.SetBit(0)
	sourceCode.SetBit(2)
	sourceCode.SetBit(3)
	sourceCode.SetBit(63)

	cluster := cs.NewCluster(sourceCode, point.GetReceptors())
	point.SetMemory(cluster)

	sourceCode.Reset()
	sourceCode.SetBit(1)
	sourceCode.SetBit(3)
	sourceCode.SetBit(5)
	sourceCode.SetBit(6)
	cluster2 := cs.NewCluster(sourceCode, point.GetReceptors())
	point.SetMemory(cluster2)

	fmt.Println(point)
}
