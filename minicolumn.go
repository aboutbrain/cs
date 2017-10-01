package cs

import (
	"fmt"
	"log"

	"github.com/golang-collections/go-datastructures/bitarray"
)

var _ = log.Printf // For debugging; delete when done.
var _ = fmt.Printf // For debugging; delete when done.

type MiniColumn struct {
	inputVector                bitarray.BitArray
	inputVectorLen             uint64
	inputLen                   int
	outputVector               bitarray.BitArray
	outputVectorLen            uint64
	learningVector             bitarray.BitArray
	cs                         *CombinatorialSpace
	clusterThreshold           int
	clusterActivationThreshold int
	epoch                      int
	memoryLimit                int
	level                      int
}

func NewMiniColumn(clusterThreshold, clusterActivationThreshold, memoryLimit int, inputVectorLen, outputVectorLen uint64) *MiniColumn {
	return &MiniColumn{
		clusterThreshold:           clusterThreshold,
		clusterActivationThreshold: clusterActivationThreshold,
		memoryLimit:                memoryLimit,
		level:                      2,
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
}

func (mc *MiniColumn) Next() {
	mc.ActivateClusters()
	mc.MakeOutVector()
	mc.ModifyClusters()
	mc.ConsolidateMemory()
}

func (mc *MiniColumn) ModifyClusters() {
	clusters := 0
	for pointId, point := range mc.cs.Points {
		deleted := 0
		for clusterId, cluster := range point.Memory {
			j := clusterId - deleted
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
				if len(clusterBits) != clusterBitsNewLen {
					if clusterBitsNewLen >= mc.clusterThreshold {
						hashOld := cluster.GetHash()
						cluster.SetNewBits(clusterBitsNew)
						hashNew := cluster.GetHash()
						if _, ok := mc.cs.OutHashSet[pointId][hashNew]; !ok {
							mc.cs.RemoveHash(pointId, hashOld)
							mc.cs.SetHash(pointId, hashNew)
						} else {
							mc.cs.DeleteCluster(&point, j)
							deleted++
							continue
						}
					} else {
						mc.cs.DeleteCluster(&point, j)
						deleted++
						continue
					}
				}
			}
			point.Memory[j] = cluster
		}
		mc.cs.Points[pointId] = point
	}
	//fmt.Printf("\nClustersCount: %d\n", clusters)
}

func (mc *MiniColumn) ConsolidateMemory() {
	for pointId, point := range mc.cs.Points {
		deleted := 0
		for clusterId, cluster := range point.Memory {
			j := clusterId - deleted
			clusterAge := mc.cs.InternalTime - cluster.startTime
			if clusterAge > 20 {
				errorFull := float32(cluster.ErrorFullCounter) / float32(cluster.ActivationFullCounter)
				if errorFull > 0.05 && cluster.ActivationFullCounter > 10 {
					mc.cs.DeleteCluster(&point, j)
					deleted++
					continue
				}
				errorPartial := float32(cluster.ErrorPartialCounter) / float32(cluster.ActivationPartialCounter)
				if errorPartial > 0.3 && cluster.ActivationPartialCounter > 10 {
					mc.cs.DeleteCluster(&point, j)
					deleted++
					continue
				}
			}
			switch {
			case cluster.Status == ClusterTmp && cluster.LearnCounter > 6:
				cluster.Status = ClusterPermanent1
			case cluster.Status == ClusterPermanent1 && cluster.LearnCounter > 16:
				cluster.Status = ClusterPermanent2
				mc.cs.clustersPermanent++
			}
			point.Memory[j] = cluster
		}
		mc.cs.Points[pointId] = point
	}
}

func (mc *MiniColumn) ActivateClusters() {
	stat := make(map[int]int)
	for pointId, point := range mc.cs.Points {
		pointPotential := 0
		clusters := 0
		for clusterId, cluster := range point.Memory {
			clusters++
			potential, inputBits := cluster.GetCurrentPotential(mc.inputVector)
			inputSize := cluster.GetInputSize()
			active, _ := mc.learningVector.GetBit(uint64(point.GetOutputBit()))
			if potential == inputSize {
				cluster.SetActivationStatus(ClusterStatusFull)
				cluster.ActivationFullCounter++
				//cluster.LearnCounterIncrease()
				if !active {
					cluster.ErrorFullCounter++
				}
				if cluster.Status == ClusterPermanent2 {
					pointPotential += potential - mc.clusterActivationThreshold + 1
				}
			} else if potential >= mc.clusterActivationThreshold {
				cluster.SetActivationStatus(ClusterStatePartial)
				cluster.ActivationPartialCounter++
				if !active {
					cluster.ErrorPartialCounter++
				} else {
					cluster.SetHistory(inputBits, active)
				}
				cluster.LearnCounterIncrease()
			}
			point.Memory[clusterId] = cluster
		}
		stat[clusters]++
		//fmt.Printf("PointId: %d, Clusters: %d\n", pointId, clusters)
		point.SetPotential(pointPotential)
		mc.cs.Points[pointId] = point
	}
	for i, v := range stat {
		fmt.Printf("Clusters: %d, Points: %d\n", i, v)
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

		outputBit := point.GetOutputBit()
		active, _ := mc.learningVector.GetBit(uint64(outputBit))

		receptorsActiveLen := len(receptorsActiveCount.ToNums())
		memorySize := len(point.Memory)

		if receptorsActiveLen >= mc.clusterThreshold && active && memorySize < mc.memoryLimit {
			cluster := NewCluster(receptorsActiveCount, mc.inputVectorLen)
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
