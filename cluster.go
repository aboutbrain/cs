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
	ClusterStatusFull
	ClusterStatePartial
)

type InputBits []uint64
type OutputBits []uint64
type History struct {
	InputBits
	OutputBits
}

type Cluster struct {
	Status                        int
	inputBitSet                   bitarray.BitArray
	targetBitSet                  bitarray.BitArray
	ActivationState               int
	potential                     int
	ActivationStateFullCounter    int
	ActivationStatePartialCounter int
	ErrorCompleteCounter          int
	ErrorPartialCounter           int
	Weights                       map[int]float32
	HistoryMemory                 []History
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

func (c *Cluster) MainComponent() {

}
