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
			clusters++
			cluster.Learn(mc.inputVector, mc.learningVector)
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
			if cluster.rHigh < 0.7 || cluster.Status == ClusterDeleting {
				point.DeleteCluster(j)
				deleted++
				mc.cs.clustersTotal--
			}
		}
	}
}

func (mc *MiniColumn) activateClustersInput() {
	stat := make(map[int]int)
	//clustersFullyActivated := 0
	//clustersPartialActivated := 0
	for pointId := range mc.cs.Points {
		point := &mc.cs.Points[pointId]
		for clusterId := range point.Memory {
			cluster := &point.Memory[clusterId]
			cluster.CalculatingInputCoincidence(mc.inputVector)
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
			for _, cluster := range point.Memory {
				if cluster.q > lowP {
					if result, _ := cluster.targetBitSet.GetBit(uint64(i)); result {
						potential += cluster.clusterResultLength
					} else {
						potential -= cluster.clusterResultLength
					}
				}
			}
		}
		if potential > 30 {
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

		inputMin := int(math.Sqrt(float64(mc.inputLen)))
		outputMin := int(math.Sqrt(float64(mc.outputLen)))

		if receptorsActiveLen >= inputMin+1 && outputsActiveLen >= outputMin+1 && memorySize < mc.memoryLimit {
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
