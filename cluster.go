package cs

import (
	"strconv"

	"github.com/golang-collections/go-datastructures/bitarray"
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

type InputBits []uint64
type OutputBits []uint64
type History struct {
	InputBits
	OutputBits
}

type Cluster struct {
	Status                   int
	startTime                int
	inputBitSet              bitarray.BitArray
	targetBitSet             bitarray.BitArray
	ActivationState          int
	potential                int
	ActivationFullCounter    int
	ActivationPartialCounter int
	ErrorFullCounter         int
	ErrorPartialCounter      int
	Weights                  map[int]float32
	HistoryMemory            []History
	LearnCounter             int
}

func NewCluster(inputBitSet, targetBitSet bitarray.BitArray) *Cluster {
	c := &Cluster{
		Status:          ClusterTmp,
		ActivationState: ClusterStateNon,
		Weights:         make(map[int]float32),
	}
	c.inputBitSet = inputBitSet
	c.targetBitSet = targetBitSet
	return c
}

func (c *Cluster) GetCurrentPotential(inputVector bitarray.BitArray) (int, InputBits) {
	inputBits := inputVector.And(c.inputBitSet).ToNums()
	c.potential = len(inputBits)
	return c.potential, inputBits
}

func (c *Cluster) SetHistory(inputBits InputBits, outputBits OutputBits) {
	c.HistoryMemory = append(c.HistoryMemory, History{InputBits: inputBits, OutputBits: outputBits})
}

func (c *Cluster) SetStatus(status int) {
	c.Status = status
}

func (c *Cluster) SetActivationStatus(status int) {
	c.ActivationState = status
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

func (c *Cluster) SetNewBits(nums []uint64) {
	a := bitarray.NewBitArray(256)
	for _, num := range nums {
		a.SetBit(num)
	}
	c.inputBitSet = a
}

func (c *Cluster) BitActivationStatistic() ([]float32, []uint64) {
	var max float32 = 0
	var a int = 0

	activeBits := c.inputBitSet.ToNums()
	clusterLength := len(activeBits)
	f := make([]float32, clusterLength)
	nu := 1 / float32(clusterLength)

	for i := range f {
		f[i] = 1.0
	}

	for j := 0; j < 2; j++ {
		for _, v := range c.HistoryMemory {
			a = 0

			for l, n := range activeBits {
				if inArray(int(n), v.InputBits) {
					a += int(f[l])
				}
			}

			for l, n := range activeBits {
				if inArray(int(n), v.InputBits) {
					fl := float32(a) * nu
					f[l] += fl
				}
			}

			for _, e := range f {
				if e > max {
					max = e
				}
			}

			for i := range f {
				f[i] = f[i] / max
			}
		}
		nu = nu * 0.8
	}
	return f, activeBits
}
