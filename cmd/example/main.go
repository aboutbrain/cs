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
	CombinatorialSpaceSize     = 10000
	BitsPerPoint               = 8
	ClusterThreshold           = 6
	ClusterActivationThreshold = 4
	CharacterBits              = 8
	PointMemoryCapacity        = 10
	PointContextCapacity       = 10
)

func main() {
	rand.Seed(time.Now().Unix())
	comSpace := cs.NewCombinatorialSpace(CombinatorialSpaceSize, BitsPerPoint, OutputVectorSize)
	point := comSpace.Points[0]

	sourceCode := bitarray.NewBitArray(InputVectorSize)
	sourceCode.SetBit(0)
	sourceCode.SetBit(2)
	sourceCode.SetBit(3)
	sourceCode.SetBit(63)

	for i, p := range comSpace.Points  {
		cluster := cs.NewCluster(sourceCode, p.GetReceptors())
		hash := cluster.GetHash()
		if comSpace.CheckOutHashSet(p.OutBit, hash){
			comSpace.SetHash(p.OutBit, hash)
			p.SetMemory(cluster)
			comSpace.Points[i] = p
		}
	}

	fmt.Println(point)
}
