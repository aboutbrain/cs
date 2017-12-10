package cs

import (
	"fmt"

	"github.com/aboutbrain/cs/bitarray"
)

func (mc *MiniColumn) CreateOutput() (bitarray.BitArray, int) {
	mc.activateClustersInput()
	clusters := mc.makeOutVector()
	return mc.outputVector, clusters
}

func (mc *MiniColumn) LearnUnsupervised(inputVector, outputVector bitarray.BitArray) {
	mc.SetInputVector(inputVector)
	mc.SetLearningVector(mc.createLearningVector(outputVector))
	result := mc.outputVector.Equals(mc.learningVector)
	if !result {
		mc.addNewClusters()
	}
}

func (mc *MiniColumn) createLearningVector(outputVector bitarray.BitArray) bitarray.BitArray {
	learningVector := bitarray.NewBitArray(mc.outputVectorLen)
	nums := outputVector.ToNums()
	for _, bitNum := range nums {
		learningVector.SetBit(bitNum)
	}
	l := len(learningVector.ToNums())
	delta := 8 - l
	for i := delta; i > 0; i-- {
	rnd:
		num := Random(0, int(mc.outputVectorLen))
		if bit, _ := learningVector.GetBit(uint64(num)); !bit {
			learningVector.SetBit(uint64(num))
		} else {
			goto rnd
		}
	}
	str, num := BitArrayToString(learningVector, int(mc.outputVectorLen))
	fmt.Printf("LearningVector:  %s, %d\n", str, num)
	return learningVector
}

func (mc *MiniColumn) EncourageClusters() {

}

func (mc *MiniColumn) BlameClusters() {

	//points := mc.cs.outBitToPointsMap[bitId]

}
