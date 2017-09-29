package cs

import (
	"fmt"

	"github.com/golang-collections/go-datastructures/bitarray"
)

var _ = fmt.Printf // For debugging; delete when done.

type MiniColumn struct {
	inputLen         int
	inputVector      bitarray.BitArray
	outputVector     bitarray.BitArray
	outputLen        int
	learningVector   bitarray.BitArray
	cs               *CombinatorialSpace
	clusterThreshold int
	epoch            int
	memoryLimit      int
	level            int
}

func NewMiniColumn(clusterThreshold, memoryLimit int) *MiniColumn {
	return &MiniColumn{
		clusterThreshold: clusterThreshold,
		memoryLimit:      memoryLimit,
		level:            50,
		outputVector:     bitarray.NewBitArray(256),
	}
}

func (mc *MiniColumn) SetCombinatorialSpace(cs *CombinatorialSpace) {
	mc.cs = cs
}

func (mc *MiniColumn) SetInputVector(inputVector bitarray.BitArray) {
	mc.inputVector = inputVector
	mc.inputLen = len(mc.inputVector.ToNums())
}

func (mc *MiniColumn) SetLearningVector(learningVector bitarray.BitArray) {
	mc.learningVector = learningVector
	mc.outputLen = len(mc.learningVector.ToNums())
}

func (mc *MiniColumn) Next() {
	mc.ActivateClusters()
	mc.MakeOutVector()
	mc.ModifyClusters()
	mc.ConsolidateMemory()
}

func (mc *MiniColumn) ModifyClusters() {

}

func (mc *MiniColumn) ConsolidateMemory() {

}

func (mc *MiniColumn) ActivateClusters() {
	for i, point := range mc.cs.Points {
		pointPotential := 0
		for j, cluster := range point.Memory {
			potential, inputBits := cluster.GetCurrentPotential(mc.inputVector)
			inputSize := cluster.GetInputSize()
			if potential == inputSize {
				cluster.SetActivationStatus(ClusterStatusFull)
				cluster.ActivationStateFullCounter++
				if mc.learningVector.Intersects(cluster.targetBitSet) {
					cluster.ErrorCompleteCounter++
				}
				if cluster.Status == ClusterPermanent2 {
					pointPotential += potential - 3 + 1
				}
			} else if potential >= 3 {
				cluster.SetActivationStatus(ClusterStatePartial)
				cluster.ActivationStatePartialCounter++
				outBits := mc.learningVector.And(cluster.targetBitSet).ToNums()
				if len(outBits) == 0 {
					cluster.ErrorPartialCounter++
				} else {
					cluster.SetHistory(inputBits, outBits)
				}
			}
			point.Memory[j] = cluster
		}
		point.SetPotential(pointPotential)
		mc.cs.Points[i] = point
	}
}

func (mc *MiniColumn) MakeOutVector() {
	for i, currentOutBitPointsMap := range mc.cs.outBitToPointsMap {
		p := 0
		for _, pointId := range currentOutBitPointsMap {
			p += mc.cs.Points[pointId].GetPotential()
		}
		if p >= mc.level {
			mc.outputVector.SetBit(uint64(i))
		}
	}
}

func (mc *MiniColumn) OutVector() bitarray.BitArray {
	return mc.outputVector
}

func (mc *MiniColumn) AddNewClusters() {
	for pointId, point := range mc.cs.Points {

		receptors := point.GetReceptors()
		receptorsActiveCount := mc.inputVector.And(receptors)

		outputs := point.GetOutputs()
		outputsActiveCount := mc.learningVector.And(outputs)

		receptorsActiveLen := len(receptorsActiveCount.ToNums())
		outputsActiveLen := len(outputsActiveCount.ToNums())
		memorySize := len(point.Memory)

		if receptorsActiveLen > mc.inputLen/3 && outputsActiveLen > mc.outputLen/3 && memorySize < mc.memoryLimit {
			cluster := NewCluster(receptorsActiveCount, outputsActiveCount)
			hash := cluster.GetHash()
			//fmt.Println(hash)
			if mc.cs.CheckOutHashSet(pointId, hash) {
				point.SetMemory(cluster)
				mc.cs.Points[pointId] = point
				mc.cs.SetHash(pointId, hash)
			}
		}
	}
	mc.epoch++
}
