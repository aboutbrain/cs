package cs

import (
	"strconv"

	"github.com/golang-collections/go-datastructures/bitarray"
)

const (
	ClusterTmp = iota
	ClusterPermanent1
	ClusterPremanent2
	ClusterDeleting
)

type Cluster struct {
	Status            int
	inputBitSet       bitarray.BitArray
	targetBitSet      bitarray.BitArray
	potential         int
	ActivationCounter int
	ErrorCounter      int
}

func NewCluster(inputVector, inputReceptors, targetVector, targetOutputs bitarray.BitArray) *Cluster {
	c := &Cluster{
		Status: ClusterTmp,
	}
	c.inputBitSet = inputReceptors.And(inputVector)
	c.targetBitSet = targetOutputs.And(targetVector)
	return c
}

func (c *Cluster) SetCurrentPotential(targetVector bitarray.BitArray) {
	//c.potential = targetVector
}

func (c *Cluster) GetSize() int {
	return len(c.inputBitSet.ToNums())
}

func (c *Cluster) GetHash() string {
	nums := c.inputBitSet.ToNums()
	hash := ""
	for _, v := range nums {
		hash += "." + strconv.Itoa(int(v))
	}
	return hash
}
