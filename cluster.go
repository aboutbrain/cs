package cs

import "github.com/golang-collections/go-datastructures/bitarray"

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
