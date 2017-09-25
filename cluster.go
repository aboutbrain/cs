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
	bitSet            bitarray.BitArray
	potential         int
	ActivationCounter int
	ErrorCounter      int
}

func NewCluster(input, point bitarray.BitArray) *Cluster {
	c := &Cluster{
		Status: ClusterTmp,
	}
	c.bitSet = point.And(input)
	return c
}

func (c *Cluster) SetCurrentPotential(targetVector bitarray.BitArray) {
	c.potential = i
}

func (c *Cluster) GetSize() int {
	return len(c.bitSet.ToNums())
}

func (c *Cluster) GetHash() string {
	nums := c.bitSet.ToNums()
	hash := ""
	for _, v := range nums {
		hash += "." + strconv.Itoa(int(v))
	}
	return hash
}
