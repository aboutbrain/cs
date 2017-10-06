package cs

import (
	"fmt"
	"log"

	"github.com/aboutbrain/cs/bitarray"
	"strconv"
)

var _ = log.Printf // For debugging; delete when done.
var _ = fmt.Printf // For debugging; delete when done.

type MiniColumn struct {
	inputVector                bitarray.BitArray
	inputVectorLen             int
	inputLen                   int
	outputVector               bitarray.BitArray
	outputVectorLen            int
	learningVector             bitarray.BitArray
	cs                         *CombinatorialSpace
	clusterThreshold           int
	clusterActivationThreshold int
	epoch                      int
	memoryLimit                int
	level                      int
}

func NewMiniColumn(clusterThreshold, clusterActivationThreshold, memoryLimit int, inputVectorLen, outputVectorLen int) *MiniColumn {
	return &MiniColumn{
		clusterThreshold:           clusterThreshold,
		clusterActivationThreshold: clusterActivationThreshold,
		memoryLimit:                memoryLimit,
		level:                      2,
		inputVectorLen:             inputVectorLen,
		outputVectorLen:            outputVectorLen,
		outputVector:               bitarray.NewBitArray(uint64(outputVectorLen)),
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

func (mc *MiniColumn) Calculate() bitarray.BitArray {
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
	mc.activateClusters()
	//mc.inputVector = inputVectorSave
	mc.makeOutVector()
	mc.modifyClusters()
	mc.consolidateMemory()
	//mc.checkClusters("Calculate")
	return mc.outputVector
}

func (mc *MiniColumn) checkClusters(method string){
	clusters := mc.clustersCount()
	hashes := mc.hashCount()
	if clusters != hashes {
		panic("Не совпадают: " + method + ", clusters:" + strconv.Itoa(clusters) + ", hashes:" + strconv.Itoa(hashes))
	}
	fmt.Printf("\nClustersCount: %d, HashesCount: %d after deleting\n", clusters, hashes)
}

func (mc *MiniColumn) clustersCount() int {
	clusters := 0
	for _, point := range mc.cs.Points {
		for range point.Memory {
			clusters++
		}
	}
	//fmt.Printf("\nClustersCount after deleting: %d\n", clusters)
	return clusters
}

func (mc *MiniColumn) hashCount() int {
	hashCount := 0
	for _, pointHash := range mc.cs.OutHashSet {
		for range pointHash {
			hashCount++
		}
	}
	//fmt.Printf("HashesCount after deleting: %d\n\n", hashCount)
	return hashCount
}

func (mc *MiniColumn) Learn(day bool) {
	if day {
		mc.addNewClusters()
		//mc.activateClusters()
	}
}

func (mc *MiniColumn) modifyClusters() {
	//log.Println("Вошли в modifyClusters!!")
	//mc.checkClusters("До modifyClusters")
	clusters := 0
	for pointId := range mc.cs.Points {
		point := &mc.cs.Points[pointId]
		deleted := 0
		for clusterId := range point.Memory {
			j := clusterId - deleted
			cluster := &point.Memory[j]
			clusters++
			if cluster.Status != ClusterPermanent2 && (cluster.LearnCounter == 4 || cluster.LearnCounter == 16) {
				f := cluster.BitActivationStatistic()
				//f := cluster.Weights
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
						//cluster.created = "уменьшили - hashOld: " + hashOld + ", hashNew: " + hashNew
						if !mc.cs.CheckOutHashSet(point.OutputBit, hashNew) {
							t, ok := mc.cs.OutHashSet[point.OutputBit][hashNew]
							t, ok = mc.cs.OutHashSet[point.OutputBit][hashOld]
							mc.cs.RemoveHash(point.OutputBit, hashOld)
							mc.cs.SetHash(point.OutputBit, hashNew)
							t, ok = mc.cs.OutHashSet[point.OutputBit][hashOld]
							t, ok = mc.cs.OutHashSet[point.OutputBit][hashNew]
							_ = t
							_ = ok

							/*w := make(map[int]float32)
							for _, v := range clusterBitsNew {
								w[int(v)] = f[int(v)]
							}
							cluster.Weights = w*/
							cluster.clusterLength = clusterBitsNewLen
							//point.Memory[j] = cluster
							//mc.cs.Points[pointId] = point
							//log.Println("Поменяли по обрезке!, pointId:" + strconv.Itoa(point.id) + ", clusterId: " + strconv.Itoa(j))
							//mc.checkClusters("modifyClusters в точке: " + strconv.Itoa(point.id) + ", clusterId: " + strconv.Itoa(j))
						} else {
							mc.cs.DeleteCluster(point, j, false)
							mc.cs.RemoveHash(point.OutputBit, hashOld)
							//mc.cs.Points[pointId] = point
							deleted++
							//log.Println("удалили по Хэшу активации!, pointId:" + strconv.Itoa(point.id) + ", clusterId: " + strconv.Itoa(j))
							//mc.checkClusters("modifyClusters в точке: " + strconv.Itoa(point.id) + ", clusterId: " + strconv.Itoa(j))
							continue
						}
					} else {
						mc.cs.DeleteCluster(point, j, true)
						//mc.cs.Points[pointId] = point
						deleted++
						//log.Println("удалили по Длине активации!, pointId:" + strconv.Itoa(point.id) + ", clusterId: " + strconv.Itoa(j))
						//mc.checkClusters("modifyClusters в точке: " + strconv.Itoa(point.id) + ", clusterId: " + strconv.Itoa(j))
						continue
					}
				}
			}
			//mc.checkClusters("modifyClusters в точке: " + strconv.Itoa(point.id))
		}
		//mc.cs.Points[pointId] = point
		//mc.checkClusters("modifyClusters в точке: " + strconv.Itoa(point.id))
	}
	//fmt.Printf("\nClustersCount after deleting: %d\n", clusters)
	//mc.checkClusters("После modifyClusters")
	//log.Println("Вышли из modifyClusters!!")
}

func (mc *MiniColumn) consolidateMemory() {
	//log.Println("Вошли в consolidateMemory!!")
	//mc.checkClusters("До consolidateMemory")
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
					//log.Println("Удалили по ActivationFullCounter!, pointId:" + strconv.Itoa(point.id) + ", clusterId: " + strconv.Itoa(j))
					//mc.checkClusters("ActivationFullCounter в точке: " + strconv.Itoa(point.id) + ", clusterId: " + strconv.Itoa(j))
					deleted++
					//log.Println("удалили по Полной активации!")
					continue
				} else if cluster.Status == ClusterTmp {
					cluster.Status = ClusterPermanent2
					mc.cs.clustersPermanent2++
				}
			}

			if cluster.ActivationPartialCounter > 5 {
				errorPartial := float32(cluster.ErrorPartialCounter) / float32(cluster.ActivationPartialCounter)
				if errorPartial > 0.3 {
					mc.cs.DeleteCluster(point, j, true)
					//log.Println("Удалили по ActivationPartialCounter!, pointId:" + strconv.Itoa(point.id) + ", clusterId: " + strconv.Itoa(j))
					//mc.checkClusters("ActivationPartialCounter в точке: " + strconv.Itoa(point.id) + ", clusterId: " + strconv.Itoa(j))
					deleted++
					//log.Println("удалили по частичной активации!")
					continue
				}
			}

			switch {
			case cluster.Status == ClusterTmp && (cluster.LearnCounter > 6):
				cluster.Status = ClusterPermanent1
				mc.cs.clustersPermanent1++
			case cluster.Status == ClusterPermanent1 && (cluster.LearnCounter > 16):
				cluster.Status = ClusterPermanent2
				mc.cs.clustersPermanent2++
			}
			//point.Memory[j] = cluster
		}
		//mc.cs.Points[pointId] = point
	}
	//mc.checkClusters("После consolidateMemory")
	//log.Println("Вышли из consolidateMemory!!")
}

func (mc *MiniColumn) activateClusters() {
	//log.Println("Вошли в activateClusters!!")
	stat := make(map[int]int)
	clustersTotal := 0
	clustersFullyActivated := 0
	clustersPartialActivated := 0
	for pointId, point := range mc.cs.Points {
		point.potential = 0
		pointPotential := 0
		clusters := 0
		for clusterId, cluster := range point.Memory {
			clusters++
			clustersTotal++
			result := cluster.inputBitSet.And(mc.inputVector)
			nums := result.ToNums()
			cluster.potential = len(nums)
			inputSize := cluster.GetInputSize()
			active, _ := mc.learningVector.GetBit(uint64(point.GetOutputBit()))
			if cluster.potential == inputSize {
				cluster.SetActivationStatus(ClusterStatusFull)
				cluster.ActivationFullCounter++
				clustersFullyActivated++
				if !active {
					cluster.ErrorFullCounter++
				}
				if cluster.Status == ClusterPermanent2 {
					pointPotential += cluster.potential - mc.clusterActivationThreshold + 1
				}
			} else if int(cluster.potential) >= mc.clusterActivationThreshold {
				clustersPartialActivated++
				cluster.SetActivationStatus(ClusterStatePartial)
				cluster.ActivationPartialCounter++
				if !active {
					cluster.ErrorPartialCounter++
				} else {
					if cluster.Status != ClusterPermanent2 {
						//cluster.BitStatisticNew(nums)

						cluster.SetHistory(nums, active)
					}
				}
				if cluster.Status != ClusterPermanent2 {
					cluster.LearnCounterIncrease()
				}
			}
			point.Memory[clusterId] = cluster
		}
		stat[clusters]++
		point.SetPotential(int(pointPotential))
		mc.cs.Points[pointId] = point
	}
	fmt.Printf("ActivatedClusters: Fully: %d, Partly: %d\n", clustersFullyActivated, clustersPartialActivated)
	fmt.Printf("\nClustersTotal before deleting: %d\n", clustersTotal)
	/*for i, v := range stat {
		fmt.Printf("Clusters: %d, Points: %d\n", i, v)
	}*/
	//log.Println("Вышли из activateClusters!!")
}

func (mc *MiniColumn) makeOutVector() {
	mc.outputVector.Reset()
	for i, currentOutBitPointsMap := range mc.cs.outBitToPointsMap {
		p := 0
		for _, pointId := range currentOutBitPointsMap {
			p += mc.cs.Points[pointId].GetPotential()
		}
		if p >= int(mc.level) {
			mc.outputVector.SetBit(uint64(i))
		}
	}
}

func (mc *MiniColumn) OutVector() bitarray.BitArray {
	return mc.outputVector
}

func (mc *MiniColumn) addNewClusters() {
	//log.Println("Вошли в addNewClusters!!")
	for pointId, point := range mc.cs.Points {

		receptors := point.GetReceptors()
		activeReceptors := mc.inputVector.And(receptors)

		outputBit := point.GetOutputBit()
		active, _ := mc.learningVector.GetBit(uint64(outputBit))

		receptorsActiveLen := len(activeReceptors.ToNums())
		memorySize := len(point.Memory)

		if receptorsActiveLen >= mc.clusterThreshold && active && memorySize < mc.memoryLimit {
			cluster := NewCluster(activeReceptors, mc.inputVectorLen)
			cluster.startTime = mc.cs.InternalTime
			hash := cluster.GetHash()
			if !mc.cs.CheckOutHashSet(outputBit, hash) {
				point.SetMemory(cluster)
				mc.cs.Points[pointId] = point
				mc.cs.SetHash(outputBit, hash)
				mc.cs.clustersTotal++
				//log.Println("Добавили новый!")
			}
		}
	}
	mc.epoch++
	//log.Println("Вышли из addNewClusters!!")
}
