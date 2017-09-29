package cs

import (
	"fmt"
	"log"

	"github.com/golang-collections/go-datastructures/bitarray"
)

var _ = log.Printf
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
	for pointId, point := range mc.cs.Points {
		for clusterId, cluster := range point.Memory {
			cluster.LearnCounter++
			//log.Printf("learnCounter: %d, Test: %d", cluster.LearnCounter, cluster.startTime)
			if cluster.Status != ClusterPermanent2 && (cluster.LearnCounter == 4 || cluster.LearnCounter == 16) {
				f, clusterBits := cluster.BitActivationStatistic()
				clusterBitsNew := []uint64{}
				for k, bitNumber := range clusterBits {
					if f[k] > 0.75 {
						clusterBitsNew = append(clusterBitsNew, bitNumber)
					}
				}
				clusterBitsNewLen := len(clusterBitsNew)
				if len(clusterBits) != clusterBitsNewLen {
					if clusterBitsNewLen > 3 {
						cluster.SetNewBits(clusterBitsNew)
						hash := cluster.GetHash()
						delete(mc.cs.OutHashSet[pointId], hash)
						mc.cs.SetHash(pointId, hash)
					} else {
						mc.cs.DeleteCluster(pointId, clusterId)
					}
					mc.cs.DeleteCluster(pointId, clusterId)
				}
			}
			point.Memory[clusterId] = cluster
		}
		mc.cs.Points[pointId] = point
	}
}

func (mc *MiniColumn) ConsolidateMemory() {
	for pointId, point := range mc.cs.Points {
		for clusterId, cluster := range point.Memory {
			if mc.cs.InternalTime-cluster.startTime > 100 {
				if float32(cluster.ErrorFullCounter)/float32(cluster.ActivationFullCounter) > 0.05 {
					mc.cs.DeleteCluster(pointId, clusterId)
					continue
				}
				if float32(cluster.ErrorPartialCounter)/float32(cluster.ActivationPartialCounter) > 0.3 {
					mc.cs.DeleteCluster(pointId, clusterId)
					continue
				}
			}
			switch {
			case cluster.Status == ClusterTmp && cluster.ActivationState == ClusterStatusFull:
				cluster.Status = ClusterPermanent2
				mc.cs.clustersPermanent++
			case cluster.Status == ClusterTmp && cluster.LearnCounter > 6:
				cluster.Status = ClusterPermanent1
			case cluster.Status == ClusterPermanent1 && cluster.LearnCounter > 16:
				cluster.Status = ClusterPermanent2
				mc.cs.clustersPermanent++
			}
			//log.Printf("---learnCounter: %d, Test: %d", cluster.LearnCounter, cluster.startTime)
			point.Memory[clusterId] = cluster
		}
		mc.cs.Points[pointId] = point
	}
}

func (mc *MiniColumn) ActivateClusters() {
	for pointId, point := range mc.cs.Points {
		pointPotential := 0
		for clusterId, cluster := range point.Memory {
			potential, inputBits := cluster.GetCurrentPotential(mc.inputVector)
			inputSize := cluster.GetInputSize()
			if potential == inputSize {
				cluster.SetActivationStatus(ClusterStatusFull)
				cluster.ActivationFullCounter++
				if !mc.learningVector.Intersects(cluster.targetBitSet) {
					cluster.ErrorFullCounter++
				}
				if cluster.Status == ClusterPermanent2 {
					pointPotential += potential - 3 + 1
				}
			} else if potential >= 3 {
				cluster.SetActivationStatus(ClusterStatePartial)
				cluster.ActivationPartialCounter++
				outBits := mc.learningVector.And(cluster.targetBitSet).ToNums()
				if len(outBits) == 0 {
					cluster.ErrorPartialCounter++
				} else {
					cluster.SetHistory(inputBits, outBits)
					//cluster.LearnCounter++
				}
				mc.cs.Points[pointId].Memory[clusterId].LearnCounter++
			}
			point.Memory[clusterId] = cluster
		}
		point.SetPotential(pointPotential)
		mc.cs.Points[pointId] = point
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

		if receptorsActiveLen > mc.inputLen/3 && outputsActiveLen >= mc.outputLen/4 && memorySize < mc.memoryLimit {
			cluster := NewCluster(receptorsActiveCount, outputsActiveCount)
			cluster.startTime = mc.cs.InternalTime
			hash := cluster.GetHash()
			if mc.cs.CheckOutHashSet(pointId, hash) {
				point.SetMemory(cluster)
				mc.cs.Points[pointId] = point
				mc.cs.SetHash(pointId, hash)
			}
		}
	}
	mc.epoch++
}
