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
		level:                      1,
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
	/*inputVectorSave := bitarray.NewBitArray(mc.inputVectorLen - 1)
	inputVectorSave = inputVectorSave.Or(mc.inputVector)
	// Add input noise
	for k := 0; k < 1; k++ {
		bitNum := Random(0, int(mc.inputVectorLen) - 1)
		mc.inputVector.ClearBit(uint64(bitNum))
	}
	if inputVectorSave.Equals(mc.inputVector) {
		fmt.Println("\033[33mШУМ!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\033[0m")
	}*/
	mc.ActivateClusters()
	//mc.inputVector = inputVectorSave
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
				f := cluster.BitActivationStatistic()
				clusterBitsNew := []uint64{}
				for i := range f {
					if f[i] > 0.75 {
						clusterBitsNew = append(clusterBitsNew, uint64(i))
					}
				}
				clusterBitsNewLen := len(clusterBitsNew)
				if len(f) != clusterBitsNewLen {
					if clusterBitsNewLen >= mc.clusterActivationThreshold {
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
			totalCounter := cluster.ActivationFullCounter + cluster.ActivationPartialCounter

			if cluster.ActivationFullCounter > 25 {
				errorFull := float32(cluster.ErrorFullCounter) / float32(cluster.ActivationFullCounter)
				if errorFull > 0.05 {
					mc.cs.DeleteCluster(&point, j)
					deleted++
					continue
				}
			}

			if cluster.ActivationPartialCounter > 5 {
				errorPartial := float32(cluster.ErrorPartialCounter) / float32(cluster.ActivationPartialCounter)
				if errorPartial > 0.3 {
					mc.cs.DeleteCluster(&point, j)
					deleted++
					continue
				}
			}

			switch {
			/*case cluster.Status == ClusterTmp && clusterAge > 100 && errorFull == 0: {
				cluster.Status = ClusterPermanent1
			}
			case cluster.Status == ClusterPermanent1 && clusterAge > 100 && errorFull == 0:
				cluster.Status = ClusterPermanent2
				mc.cs.clustersPermanent++*/
			case cluster.Status == ClusterTmp && (cluster.LearnCounter > 6 || totalCounter > 25):
				cluster.Status = ClusterPermanent1
				mc.cs.clustersPermanent1++
			case cluster.Status == ClusterPermanent1 && (cluster.LearnCounter > 16  || totalCounter > 50):
				cluster.Status = ClusterPermanent2
				mc.cs.clustersPermanent2++
			}
			point.Memory[j] = cluster
		}
		mc.cs.Points[pointId] = point
	}
}

func (mc *MiniColumn) ActivateClusters() {
	stat := make(map[int]int)
	clustersFullyActivated := 0
	clustersPartialActivated := 0
	for pointId, point := range mc.cs.Points {
		point.potential = 0
		pointPotential := 0
		clusters := 0
		for clusterId, cluster := range point.Memory {
			clusters++
			//totalCounter := cluster.ActivationFullCounter + cluster.ActivationPartialCounter
			//clusterAge := mc.cs.InternalTime - cluster.startTime
			potential, inputBits := cluster.GetCurrentPotential(mc.inputVector)
			inputSize := cluster.GetInputSize()
			//historyLen := len(cluster.HistoryMemory)
			active, _ := mc.learningVector.GetBit(uint64(point.GetOutputBit()))
			if potential == inputSize {
				cluster.SetActivationStatus(ClusterStatusFull)
				cluster.ActivationFullCounter++
				clustersFullyActivated++
				if !active {
					cluster.ErrorFullCounter++
				}
				if cluster.Status == ClusterPermanent2 {
					pointPotential += potential - mc.clusterActivationThreshold + 1
				}
				/*switch {
				case cluster.LearnCounter > 5:
					{
						fmt.Printf("\033[33mClusterFull\033[0m: %d, %+v \n", cluster.LearnCounter, cluster)
					}
					case cluster.LearnCounter > 0 && cluster.LearnCounter < 7: {
						fmt.Printf("\033[33mClusterFull\033[0m:%d, age: %d, %+v \n", totalCounter, clusterAge, cluster)
					}
					case cluster.LearnCounter >= 7 && cluster.LearnCounter < 16: {
						fmt.Printf("\033[32mClusterFull\033[0m:%d, age:, %d, %+v \n", totalCounter, clusterAge, cluster)
					}
					default:
						fmt.Printf("ClusterFull:%d, age: %d, %+v \n", totalCounter, clusterAge, cluster)
				}*/
			} else if potential >= mc.clusterActivationThreshold {
				clustersPartialActivated++
				cluster.SetActivationStatus(ClusterStatePartial)
				cluster.ActivationPartialCounter++
				if !active {
					cluster.ErrorPartialCounter++
				} else {
					if cluster.Status != ClusterPermanent2 {
						cluster.SetHistory(inputBits, active)
						//cluster.LearnCounterIncrease()
					}
				}
				if cluster.Status != ClusterPermanent2 {
					//cluster.SetHistory(inputBits, active)
					cluster.LearnCounterIncrease()
				}

			}
			point.Memory[clusterId] = cluster
		}
		stat[clusters]++
		//fmt.Printf("PointId: %d, Clusters: %d\n", pointId, clusters)
		point.SetPotential(pointPotential)
		mc.cs.Points[pointId] = point
	}
	fmt.Printf("Clusters: Fully: %d, Partly: %d\n", clustersFullyActivated, clustersPartialActivated)
	/*for i, v := range stat {
		fmt.Printf("Clusters: %d, Points: %d\n", i, v)
	}*/
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
