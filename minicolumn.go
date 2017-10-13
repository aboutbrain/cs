package cs

import (
	"fmt"
	"log"

	"math"

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
	needActivate               bool
	InputText                  string
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

func (mc *MiniColumn) SetInputVector(inputVector bitarray.BitArray) int {
	mc.inputVector = inputVector
	mc.inputLen = len(mc.inputVector.ToNums())
	return mc.inputLen
}

func (mc *MiniColumn) SetLearningVector(learningVector bitarray.BitArray) int {
	mc.learningVector = learningVector
	mc.outputLen = len(mc.learningVector.ToNums())
	return mc.outputLen
}

func (mc *MiniColumn) Calculate() bitarray.BitArray {
	mc.activateClustersInput()
	mc.makeOutVector()
	return mc.outputVector
}

func (mc *MiniColumn) Learn(day bool) {
	mc.needActivate = false
	if mc.needActivate {
		mc.activateClustersInput()
	}
	mc.activateClustersOutput()
	mc.modifyClusters()
	if day {
		mc.addNewClusters()
	}
	mc.calculateCorrelation()
	mc.consolidateMemory()
}

func (mc *MiniColumn) modifyClusters() {
	clusters := 0
	for pointId := range mc.cs.Points {
		point := &mc.cs.Points[pointId]
		deleted := 0
		for clusterId := range point.Memory {
			j := clusterId - deleted
			cluster := &point.Memory[j]
			oldHash := cluster.GetHash()
			clusters++
			cluster.Learn(mc.inputVector, mc.learningVector)
			newHash := cluster.GetHash()
			if oldHash != newHash {
				if !mc.cs.CheckOutHashSet(pointId, newHash) {
					mc.cs.RemoveHash(pointId, oldHash)
					mc.cs.SetHash(pointId, newHash)
				} else {
					mc.cs.DeleteCluster(point, j, false)
					mc.cs.RemoveHash(pointId, oldHash)
					deleted++
				}
			}
		}
	}
}

func (mc *MiniColumn) calculateCorrelation() {
	for pointId := range mc.cs.Points {
		point := &mc.cs.Points[pointId]
		for clusterId := range point.Memory {
			cluster := &point.Memory[clusterId]
			cluster.CalculateCorrelation()
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
			if cluster.rHigh < 0.70 || cluster.Status == ClusterDeleting {
				mc.cs.DeleteCluster(point, j, true)
				deleted++
			}
		}
		/*if deleted > 0 {
			fmt.Printf("У точки %d удалено %d кластеров\n", pointId, deleted)
		}*/
	}
}

func (mc *MiniColumn) activateClustersInput() {
	stat := make(map[int]int)
	//clustersFullyActivated := 0
	//clustersPartialActivated := 0
	for pointId := range mc.cs.Points {
		point := &mc.cs.Points[pointId]
		point.activated = 0
		for clusterId := range point.Memory {
			cluster := &point.Memory[clusterId]
			cluster.CalculatingInputCoincidence(mc.inputVector)
			if cluster.inputCoincidence == 1 && cluster.q > 0.85 {
				point.activated++
			}
		}
	}
	for i, v := range stat {
		fmt.Printf("Clusters: %d, Points: %d\n", i, v)
	}
	//fmt.Printf("ActivatedClusters: Fully: %d, Partly: %d\n", clustersFullyActivated, clustersPartialActivated)
}

func (mc *MiniColumn) activateClustersOutput() {
	for pointId := range mc.cs.Points {
		point := &mc.cs.Points[pointId]
		for clusterId := range point.Memory {
			cluster := &point.Memory[clusterId]
			cluster.CalculatingOutputCoincidence(mc.learningVector)
		}
	}
}

func (mc *MiniColumn) makeOutVector() {
	const lowP = float32(0.85)
	mc.outputVector.Reset()
	for i, currentOutBitPointsMap := range mc.cs.outBitToPointsMap {
		potential := float32(0)
		for _, pointId := range currentOutBitPointsMap {
			point := &mc.cs.Points[pointId]
			activated := 0
			for _, cluster := range point.Memory {
				if cluster.q > lowP {
					if cluster.inputCoincidence == 1 {
						if result, _ := cluster.targetBitSet.GetBit(uint64(i)); result {
							potential += cluster.inputCoincidence
							activated++
						} else {
							if cluster.outputWeights[i] > 0 {
								potential -= cluster.outputWeights[i]
							}
						}
					}
				}
			}
			if activated > 1 {
				fmt.Printf("Точка %d активирована с потенциалом %d\n", pointId, activated)
			}
		}
		if potential >= 3 {
			mc.outputVector.SetBit(uint64(i))
		}
	}
}

func (mc *MiniColumn) OutVector() bitarray.BitArray {
	return mc.outputVector
}

func (mc *MiniColumn) addNewClusters() {
	clusters := 0
	exist := 0
	activated := 0
	for pointId, point := range mc.cs.Points {
		if mc.inputLen > 0 {
			if point.activated > 0 {
				//fmt.Printf("Точка %d уже активна с потенциалом %d!!!\n", pointId, point.activated)
				activated++
				continue
			}
			receptors := point.GetReceptors()
			activeReceptors := mc.inputVector.And(receptors)

			outputs := point.GetOutputs()
			outputsActiveCount := mc.learningVector.And(outputs)

			receptorsActiveLen := len(activeReceptors.ToNums())
			outputsActiveLen := len(outputsActiveCount.ToNums())
			memorySize := len(point.Memory)

			inputMin := int(math.Sqrt(float64(3 * mc.inputLen)))
			outputMin := int(math.Sqrt(float64(3 * mc.outputLen)))

			if receptorsActiveLen >= inputMin && outputsActiveLen >= outputMin && memorySize < mc.memoryLimit {
				cluster := NewCluster(activeReceptors, outputsActiveCount, mc.inputVectorLen)
				cluster.startTime = mc.cs.InternalTime
				hash := cluster.GetHash()
				if !mc.cs.CheckOutHashSet(pointId, hash) {
					cluster.textFragment = mc.InputText
					cluster.initHash = hash
					point.SetMemory(cluster)
					mc.cs.Points[pointId] = point
					mc.cs.SetHash(pointId, hash)
					mc.cs.clustersTotal++
					clusters++
				} else {
					exist++
				}
			}
		}
	}
	mc.epoch++
	fmt.Printf("Создано %d новых класеров, пропущено %d, активированных точек %d\n", clusters, exist, activated)
}
