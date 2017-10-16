package cs

import (
	"strconv"

	"math"

	"fmt"

	"github.com/aboutbrain/cs/bitarray"
)

const (
	ClusterNormal = iota
	ClusterDeleting
)

const (
	ClusterStateNon = iota
	ClusterStatePartial
	ClusterStatusFull
)

const Tb = float32(3)
const Tb2 = Tb * Tb

type InputBits []uint64
type OutputBits []uint64
type History struct {
	InputBits
	OutputBits
}

type Cluster struct {
	Status               int
	startTime            int
	inputBitSet          bitarray.BitArray
	targetBitSet         bitarray.BitArray
	textFragment         string
	initHash             string
	inputCoincidence     float32
	outputCoincidence    float32
	inputLen, outputLen  uint64
	clusterResultLength  float32
	clusterResultLength1 float32
	clusterTargetLength  float32
	clusterTargetLength1 float32
	rHigh                float32
	rLow                 float32
	q                    float32
	R                    float32
	// Accumulative
	nLearn        float64
	inputWeights  map[int]float32
	outputWeights map[int]float32
	sX            float32
	sY            float32
	sXY           float32
	realSXY       float32
}

func NewCluster(inputBitSet, targetBitSet bitarray.BitArray, inputLen, outputLen uint64) *Cluster {
	c := &Cluster{
		Status:        ClusterNormal,
		inputWeights:  make(map[int]float32),
		outputWeights: make(map[int]float32),
		inputLen:      inputLen,
		outputLen:     outputLen,
		nLearn:        1,
		rHigh:         1,
		rLow:          0.6,
		sXY:           1,
		realSXY:       1,
		sX:            1,
		sY:            1,
	}
	c.inputBitSet = inputBitSet
	inputNums := c.inputBitSet.ToNums()
	c.targetBitSet = targetBitSet
	outputNums := c.targetBitSet.ToNums()

	for _, n := range inputNums {
		c.inputWeights[int(n)] = 1
	}

	for _, n := range outputNums {
		c.outputWeights[int(n)] = 1
	}

	return c
}

func (c *Cluster) SetStatus(status int) {
	c.Status = status
}

func (c *Cluster) GetInputSize() int {
	return len(c.inputBitSet.ToNums())
}

func (c *Cluster) CalculatingInputCoincidence(inputVector bitarray.BitArray) {
	c.clusterResultLength = 0
	c.clusterResultLength1 = 0
	c.q = 0

	s := float32(0)
	s1 := 0

	for i := range c.inputWeights {
		bitValue := float32(0)
		if v, _ := inputVector.GetBit(uint64(i)); v {
			bitValue = 1
		} else {
			bitValue = 0
		}
		c.clusterResultLength += bitValue * c.inputWeights[i]

		s += c.inputWeights[i]

		if c.inputWeights[i] > 0 {
			c.clusterResultLength1 += bitValue
			s1++
		}
	}
	c.inputCoincidence = c.clusterResultLength / s
	c.clusterResultLength1 = c.clusterResultLength1 / float32(s1)
	c.q = c.rLow * c.inputCoincidence
}

func (c *Cluster) CalculatingOutputCoincidence(inputVector bitarray.BitArray) {
	c.clusterTargetLength = 0
	c.clusterTargetLength1 = 0

	s := float32(0)
	s1 := 0

	for i := range c.outputWeights {
		bitValue := float32(0)
		if v, _ := inputVector.GetBit(uint64(i)); v {
			bitValue = 1
		} else {
			bitValue = 0
		}
		c.clusterTargetLength += bitValue * c.outputWeights[i]

		s += c.outputWeights[i]

		if c.outputWeights[i] > 0 {
			c.clusterTargetLength1 += bitValue
			s1++
		}
	}
	c.outputCoincidence = c.clusterTargetLength / s
	c.clusterTargetLength1 = c.clusterTargetLength1 / float32(s1)
}

func (c *Cluster) CalculateCorrelation() {
	const fixSXY = float32(30)
	if c.clusterResultLength1 == 1 || c.clusterTargetLength1 == 1 {
		s := float32(c.clusterResultLength1 * c.clusterTargetLength1)
		c.sXY += s
		c.realSXY += s
		k := float32(1)
		if c.sXY > fixSXY {
			k = fixSXY / c.sXY
		}
		c.sXY *= k
		c.sX = (c.sX + c.clusterResultLength1) * k
		c.sY = (c.sY + c.clusterTargetLength1) * k
		c.R = c.sXY / float32(math.Sqrt(float64(c.sX*c.sY)))
		c.rLow, c.rHigh = c.fTR()
	}
}

func (c *Cluster) fTR() (float32, float32) {
	ftrLow := float32(0.8)
	ftrHigh := float32(1)
	s := c.sX + c.sY
	if s > 4 {
		i := c.R + Tb2/(2*s)
		f := Tb * float32(math.Sqrt(float64(c.R*(1-c.R)/s+Tb2/(4*s*s))))
		i2 := 1 + Tb2/s
		ftrLow = (i - f) / i2
		ftrHigh = (i + f) / i2
	}
	return ftrLow, ftrHigh
}

func (c *Cluster) Learn(inputVector, learningVector bitarray.BitArray) {
	const LearnDelete = 0.5

	if c.clusterResultLength < 3 || c.clusterTargetLength < 3 {
		return
	}

	s := math.Sqrt(float64(c.inputCoincidence * c.outputCoincidence))

	var max float32 = 0
	c.nLearn += s

	nu := math.Max(1/c.nLearn, 0.1)
	resultVector := inputVector.And(c.inputBitSet)
	resultNums := resultVector.ToNums()

	targetOutputVector := learningVector.And(c.targetBitSet)
	targetOutputNums := targetOutputVector.ToNums()

	activeBits := c.inputBitSet.ToNums()
	targetNums := c.targetBitSet.ToNums()

	for _, v := range activeBits {
		if InArray64(int(v), resultNums) && c.inputWeights[int(v)] != 0 {
			c.inputWeights[int(v)] += float32(s * nu)
		}
	}

	for _, e := range c.inputWeights {
		if e > max {
			max = e
		}
	}

	s1 := float32(0)
	inputLen := len(c.inputWeights)
	for i := range c.inputWeights {
		c.inputWeights[i] = c.inputWeights[i] / max
		if c.inputWeights[i] < LearnDelete {
			delete(c.inputWeights, i)
		}
		s1 += c.inputWeights[i]
	}
	if inputLen > len(c.inputWeights) {
		c.setNewInputBits()
	}

	if s1 <= 2 {
		c.Status = ClusterDeleting
		fmt.Printf("кластер - удален по длине!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
		return
	}

	for _, v := range targetNums {
		if InArray64(int(v), targetOutputNums) && c.outputWeights[int(v)] != 0 {
			c.outputWeights[int(v)] += float32(s * nu)
		}
	}
	for _, e := range c.outputWeights {
		if e > max {
			max = e
		}
	}

	outputLen := len(c.outputWeights)
	for i := range c.outputWeights {
		c.outputWeights[i] = c.outputWeights[i] / max
		if c.outputWeights[i] < LearnDelete {
			delete(c.outputWeights, i)
		}
		s1 += c.outputWeights[i]
	}
	if outputLen > len(c.outputWeights) {
		c.setNewOutputBits()
	}
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

func (c *Cluster) setNewInputBits() {
	a := bitarray.NewBitArray(c.inputLen)
	for num := range c.inputWeights {
		a.SetBit(uint64(num))
	}
	c.inputBitSet = a
}

func (c *Cluster) setNewOutputBits() {
	a := bitarray.NewBitArray(c.outputLen)
	for num := range c.outputWeights {
		a.SetBit(uint64(num))
	}
	c.targetBitSet = a
}
