package cs

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"strconv"
)

const (
	ClusterTmp = iota
	ClusterPermanent1
	ClusterPremanent2
	ClusterDeleting
)

type Cluster struct {
	Status int
	bitSet bitarray.BitArray
}

func NewCluster(input, point bitarray.BitArray) *Cluster {
	c := &Cluster{
		Status: ClusterTmp,
	}
	c.bitSet = point.And(input)
	return c
}

func (c *Cluster) GetHash() string {
	nums := c.bitSet.ToNums()
	hash := ""
	for _, v := range nums {
		hash += "." + strconv.Itoa(int(v))
	}
	return hash
}