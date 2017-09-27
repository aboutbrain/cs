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

func NewCluster(inputBitSet, targetBitSet bitarray.BitArray) *Cluster {
	c := &Cluster{
		Status: ClusterTmp,
	}
	c.inputBitSet = inputBitSet
	c.targetBitSet = targetBitSet
	return c
}

func (c *Cluster) SetCurrentPotential(inputVector, targetVector bitarray.BitArray) int {
	inputPotential := len(inputVector.And(c.inputBitSet).ToNums())
	outputPotential := len(targetVector.And(c.targetBitSet).ToNums())
	c.potential = inputPotential + outputPotential
	return c.potential
}

func (c *Cluster) GetInputSize() int {
	return len(c.inputBitSet.ToNums())
}

func (c *Cluster) GetHash() string {
	nums := c.inputBitSet.ToNums()
	hash := ""
	for _, v := range nums {
		hash += "." + strconv.Itoa(int(v))
	}
	hash += "-"
	nums = c.targetBitSet.ToNums()
	for _, v := range nums {
		hash += "." + strconv.Itoa(int(v))
	}
	return hash
}
