package cs

import (
	"fmt"
	"log"

	"github.com/aboutbrain/cs/bitarray"
)

var _ = log.Printf // For debugging; delete when done.
var _ = fmt.Printf // For debugging; delete when done.

type MiniColumn struct {
	inputVector                bitarray.BitArray
	inputVectorLen             uint64
	inputLen                   int
	outputVector               bitarray.BitArray
	outputVectorLen            uint64
	outputLen                  int
	learningVector             bitarray.BitArray
	cs                         *CombinatorialSpace
	clusterThreshold           int
	clusterActivationThreshold int
	epoch                      int
	memoryLimit                int
	level                      int
}

func NewMiniColumn(clusterThreshold, clusterActivationThreshold, memoryLimit int, inputVectorLen, outputVectorLen uint64, level int) *MiniColumn {
	return &MiniColumn{
		clusterThreshold:           clusterThreshold,
		clusterActivationThreshold: clusterActivationThreshold,
		memoryLimit:                memoryLimit,
		level:                      level,
		inputVectorLen:             inputVectorLen,
		outputVectorLen:            outputVectorLen,
		outputVector:               bitarray.NewBitArray(outputVectorLen),
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

func (mc *MiniColumn) Calculate() bitarray.BitArray {
	mc.activateClusters()
	mc.makeOutVector()
	mc.modifyClusters()
	mc.consolidateMemory()
	return mc.outputVector
}

func (mc *MiniColumn)  Learn(day bool) {
	if day {
		mc.addNewClusters()
	}
}

func (mc *MiniColumn) modifyClusters() {
	clusters := 0
	for pointId := range mc.cs.Points {
		point := &mc.cs.Points[pointId]
		deleted := 0
		for clusterId := range point.Memory {
			j := clusterId - deleted
			cluster := &point.Memory[j]
			clusters++
			if cluster.Status != ClusterPermanent2 && (cluster.LearnCounter == 4 || cluster.LearnCounter == 16) {
				f, clusterBits := cluster.BitActivationStatistic()
				clusterBitsNew := []uint64{}
				for i, bitNumber := range clusterBits {
					if f[i] > 0.75 {
						clusterBitsNew = append(clusterBitsNew, bitNumber)
					}
				}
				clusterBitsNewLen := len(clusterBitsNew)
				if len(f) != clusterBitsNewLen {
					if clusterBitsNewLen > mc.clusterActivationThreshold {
						hashOld := cluster.GetHash()
						cluster.SetNewBits(clusterBitsNew)
						hashNew := cluster.GetHash()
						if _, ok := mc.cs.OutHashSet[pointId][hashNew]; !ok {
							mc.cs.RemoveHash(pointId, hashOld)
							mc.cs.SetHash(pointId, hashNew)
							cluster.clusterLength = clusterBitsNewLen
						} else {
							mc.cs.DeleteCluster(point, j, false)
							mc.cs.RemoveHash(pointId, hashOld)
							deleted++
							continue
						}
					} else {
						mc.cs.DeleteCluster(point, j, true)
						deleted++
						continue
					}
				}
			}
		}
	}
}

func (mc *MiniColumn) consolidateMemory() {
	for pointId := range mc.cs.Points {
		point := &mc.cs.Points[pointId]
		deleted := 0
		for clusterId := range point.Memory {
			j := clusterId - deleted
			cluster := &point.Memory[j]
			if cluster.ActivationFullCounter > 20 {
				errorFull := float32(cluster.ErrorFullCounter) / float32(cluster.ActivationFullCounter)
				if errorFull > 0.05 {
					mc.cs.DeleteCluster(point, j, true)
					deleted++
					continue
				} else if cluster.Status == ClusterTmp {
					cluster.Status = ClusterPermanent2
					mc.cs.clustersPermanent2++
				}

				if cluster.ActivationPartialCounter > 5 {
					errorPartial := float32(cluster.ErrorPartialCounter) / float32(cluster.ActivationPartialCounter)
					if errorPartial > 0.3 {
						mc.cs.DeleteCluster(point, j, true)
						deleted++
						continue
					}
				}
			}

			switch {
			case cluster.Status == ClusterTmp && cluster.LearnCounter > 6:
				cluster.Status = ClusterPermanent1
				mc.cs.clustersPermanent1++
			case cluster.Status == ClusterPermanent1 && cluster.LearnCounter > 16:
				cluster.Status = ClusterPermanent2
				mc.cs.clustersPermanent2++
			}
		}
	}
}

func (mc *MiniColumn) activateClusters() {
	stat := make(map[int]int)
	clustersTotal := 0
	clustersFullyActivated := 0
	clustersPartialActivated := 0
	for pointId := range mc.cs.Points {
		point := &mc.cs.Points[pointId]
		point.potential = 0
		pointPotential := 0
		clusters := 0
		for clusterId := range point.Memory {
			cluster := &point.Memory[clusterId]
			clusters++
			clustersTotal++
			result := cluster.inputBitSet.And(mc.inputVector)
			nums := result.ToNums()
			cluster.potential = len(nums)
			//potential, inputBits := cluster.GetCurrentPotential(mc.inputVector)
			inputSize := cluster.GetInputSize()
			if cluster.potential == inputSize {
				cluster.SetActivationStatus(ClusterStatusFull)
				cluster.ActivationFullCounter++
				clustersFullyActivated++
				if !mc.learningVector.Intersects(cluster.targetBitSet) {
					cluster.ErrorFullCounter++
				}
				if cluster.Status == ClusterPermanent2 {
					pointPotential += cluster.potential - mc.clusterActivationThreshold + 1
				}
			} else if int(cluster.potential) >= mc.clusterActivationThreshold  {
				clustersPartialActivated++
				cluster.SetActivationStatus(ClusterStatePartial)
				cluster.ActivationPartialCounter++
				outBits := mc.learningVector.And(cluster.targetBitSet).ToNums()
				if len(outBits) == 0 {
					cluster.ErrorPartialCounter++
				} else {
					if cluster.Status != ClusterPermanent2 {
						cluster.SetHistory(nums, outBits)
					}
					//cluster.SetHistory(inputBits, outBits)
				}
				if cluster.Status != ClusterPermanent2 {
					cluster.LearnCounterIncrease()
				}
			}
		}
		stat[clusters]++
		point.SetPotential(pointPotential)
	}
	for i, v := range stat {
		fmt.Printf("Clusters: %d, Points: %d\n", i, v)
	}
	fmt.Printf("ActivatedClusters: Fully: %d, Partly: %d\n", clustersFullyActivated, clustersPartialActivated)
}

func (mc *MiniColumn) makeOutVector() {
	mc.outputVector.Reset()
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

func (mc *MiniColumn) addNewClusters() {
	for pointId, point := range mc.cs.Points {

		receptors := point.GetReceptors()
		activeReceptors := mc.inputVector.And(receptors)

		outputs := point.GetOutputs()
		outputsActiveCount := mc.learningVector.And(outputs)

		receptorsActiveLen := len(activeReceptors.ToNums())
		outputsActiveLen := len(outputsActiveCount.ToNums())
		memorySize := len(point.Memory)

		if receptorsActiveLen >= 3 && outputsActiveLen >= 2 && memorySize < mc.memoryLimit {
			cluster := NewCluster(activeReceptors, outputsActiveCount, mc.inputVectorLen)
			cluster.startTime = mc.cs.InternalTime
			hash := cluster.GetHash()
			if !mc.cs.CheckOutHashSet(pointId, hash) {
				point.SetMemory(cluster)
				mc.cs.Points[pointId] = point
				mc.cs.SetHash(pointId, hash)
				mc.cs.clustersTotal++
			}
		}
	}
	mc.epoch++
}
