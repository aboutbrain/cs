package cs

import (
	"strconv"

	"github.com/aboutbrain/cs/bitarray"
)

const (
	ClusterTmp = iota
	ClusterPermanent1
	ClusterPermanent2
	ClusterDeleting
)

const (
	ClusterStateNon = iota
	ClusterStatePartial
	ClusterStatusFull
)

type InputBits []uint8
type History struct {
	InputBits
	OutputBit bool
}

type Cluster struct {
	Status                   uint8
	startTime                int
	inputBitSet              bitarray.BitArray
	ActivationState          uint8
	potential                int
	ActivationFullCounter    int
	ActivationPartialCounter int
	ErrorFullCounter         int
	ErrorPartialCounter      int
	Weights                  map[int]float32
	LearnCounter             int
	inputLen                 int
	clusterLength            int
}

func NewCluster(inputBitSet bitarray.BitArray, inputLen int) *Cluster {
	w := make(map[int]float32)
	nums := inputBitSet.ToNums()
	for _, v := range nums {
		w[int(v)] = 1
	}
	return &Cluster{
		Status:          ClusterTmp,
		ActivationState: ClusterStateNon,
		Weights:         w,
		inputLen:        inputLen,
		inputBitSet:     inputBitSet,
		clusterLength:   len(nums),
	}
}

func (c *Cluster) SetStatus(status uint8) {
	c.Status = status
}

func (c *Cluster) LearnCounterIncrease() {
	c.LearnCounter++
}

func (c *Cluster) SetActivationStatus(status uint8) {
	c.ActivationState = status
}

func (c *Cluster) GetInputSize() int {
	return c.clusterLength
}

func (c *Cluster) GetHash() string {
	hash := ""
	for i := range c.Weights {
		hash += "." + strconv.Itoa(i)
	}
	return hash
}

func (c *Cluster) SetNewBits(nums []uint64) {
	a := bitarray.NewBitArray(uint64(c.inputLen))
	for _, num := range nums {
		a.SetBit(num)
	}
	c.inputBitSet = a
}

func (c *Cluster) BitStatisticNew(resultNums []uint64) {
	var max float32 = 0
	var a float32 = 0
	nu := 1 / float32(c.clusterLength)

	for j := 0; j < 1; j++ {
		a = 0

		for _, n := range resultNums {
			a += c.Weights[int(n)]
		}

		for _, n := range resultNums {
			c.Weights[int(n)] += a * nu
		}

		max = 0
		for _, e := range c.Weights {
			if e > max {
				max = e
			}
		}

		for i := range c.Weights {
			c.Weights[i] = c.Weights[i] / max
		}
	}
}
