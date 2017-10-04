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
	potential                uint8
	ActivationFullCounter    int
	ActivationPartialCounter int
	ErrorFullCounter         int
	ErrorPartialCounter      int
	/*Weights                  map[int]float32*/
	HistoryMemory []History
	LearnCounter  int
	inputLen      uint8
}

func NewCluster(inputBitSet bitarray.BitArray, inputLen uint8) *Cluster {
	return &Cluster{
		Status:          ClusterTmp,
		ActivationState: ClusterStateNon,
		/*Weights:         make(map[int]float32),*/
		inputLen:    inputLen,
		inputBitSet: inputBitSet,
		//HistoryMemory: make([]History, 0, 50),
	}
}

func (c *Cluster) GetCurrentPotential(inputVector bitarray.BitArray) (int, []uint8) {
	inputBits := inputVector.And(c.inputBitSet).ToNums()
	l := len(inputBits)
	c.potential = uint8(l)
	arr := make([]uint8, l, l)
	for i, v := range inputBits {
		arr[i] = uint8(v)
	}
	return int(c.potential), arr
}

func (c *Cluster) SetHistory(inputBits InputBits, active bool) {
	c.HistoryMemory = append(c.HistoryMemory, History{InputBits: inputBits, OutputBit: active})
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
	return len(c.inputBitSet.ToNums())
}

func (c *Cluster) GetHash() string {
	bitNums := c.inputBitSet.ToNums()
	hash := ""
	for _, v := range bitNums {
		hash += "." + strconv.Itoa(int(v))
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

func (c *Cluster) BitActivationStatistic() (map[int]float32) {
	var max float32 = 0
	//var a int = 0
	var a float32 = 0

	activeBits := c.inputBitSet.ToNums()
	clusterLength := len(activeBits)
	f := make(map[int]float32, clusterLength)
	nu := 1 / float32(clusterLength)

	for _, num := range activeBits {
		f[int(num)] = 1.0
	}

	for j := 0; j < 2; j++ {
		for _, item := range c.HistoryMemory {
			a = 0

			for _, n := range activeBits {
				if InArray8(int(n), item.InputBits) {
					//a += int(f[int(n)])
					a += f[int(n)]
				}
			}

			for _, n := range activeBits {
				if InArray8(int(n), item.InputBits) {
					fl := float32(a) * nu
					f[int(n)] += fl
				}
			}

			max = 0
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
	return f
}
