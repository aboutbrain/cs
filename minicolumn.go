package cs

import (
	"fmt"
	"log"

	"github.com/aboutbrain/cs/bitarray"
)

var _ = log.Printf // For debugging; delete when done.
var _ = fmt.Printf // For debugging; delete when done.

type votingCounters struct {
	Up           int
	Down         int
	Potential    float32
	PointCluster map[int]map[int]*Cluster
}

type MiniColumn struct {
	inputVector                bitarray.BitArray // Input vector of minicolumn
	inputVectorLen             uint64            // Length of input vector
	inputLen                   int               // Number of bit in input Vector
	outputVector               bitarray.BitArray // Output vector
	outputVectorLen            uint64            // Length of output vector
	learningVector             bitarray.BitArray
	learningLen                int                 // Number of bit in learning Vector
	cs                         *CombinatorialSpace //
	clusterThreshold           int                 // cluster threshold in point
	clusterActivationThreshold int                 // cluster activation threshold
	epoch                      int
	memoryLimit                int // Limit of clusters in point
	level                      int // output bit activation level
	needActivate               bool
	InputText                  string // input text fragment for debugging only
	votingArray                map[int]votingCounters
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
		//votingArray:                make(map[int]map[int]bool),
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
	mc.learningLen = len(mc.learningVector.ToNums())
	return mc.learningLen
}

func (mc *MiniColumn) GetLearningVector() bitarray.BitArray {
	return mc.learningVector
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
					cluster.currentHash = newHash
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
	for pointId := range mc.cs.Points {
		point := &mc.cs.Points[pointId]
		point.activated = 0
		for clusterId := range point.Memory {
			cluster := &point.Memory[clusterId]
			cluster.CalculatingInputCoincidence(mc.inputVector)
			if cluster.inputCoincidence == 1 && cluster.q > 0.8 {
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

func (mc *MiniColumn) makeOutVector() int {
	mc.votingArray = make(map[int]votingCounters)
	const lowP = float32(0.86)
	clusters := 0
	mc.outputVector = bitarray.NewBitArray(mc.outputVectorLen)
	for i, currentOutBitPointsMap := range mc.cs.outBitToPointsMap {
		potential := float32(0)
		potMap := votingCounters{PointCluster: make(map[int]map[int]*Cluster)}
		for _, pointId := range currentOutBitPointsMap {
			point := &mc.cs.Points[pointId]
			activated := 0
			clMap := make(map[int]*Cluster)
			for clusterId, cluster := range point.Memory {
				if cluster.q > lowP || cluster.goodCounter < 10 {
					if cluster.inputCoincidence == 1 {
						if result, _ := cluster.targetBitSet.GetBit(uint64(i)); result {
							potential += cluster.outputWeights[i]
							clMap[clusterId] = &point.Memory[clusterId]
							potMap.Up++
							potMap.Potential += cluster.outputWeights[i]
							activated++
							clusters++
						} else {
							if cluster.outputWeights[i] == 0 {
								potMap.Potential -= 1
								potential -= 1
								clMap[clusterId] = &point.Memory[clusterId]
								potMap.Down++
							}
						}
					}
				}
			}
			if activated > 0 {
				potMap.PointCluster[pointId] = clMap
				if activated > 1 {
					fmt.Printf("Точка %d активирована с потенциалом %d\n", pointId, activated)
				}
			}
		}
		if potential >= float32(mc.level) {
			mc.votingArray[i] = potMap
			mc.outputVector.SetBit(uint64(i))
		}
	}
	nums := mc.outputVector.ToNums()
	if len(nums) > 8 {
		for _, v := range mc.votingArray {
			fmt.Printf("Potential: %f, Up: %d, Down: %d, Result: %d\n", v.Potential, v.Up, v.Down, v.Up-v.Down)
		}
		min, max := MinMax(mc.votingArray)
		fmt.Printf("Min: %d, max: %d\n", min, max)
	}
	return clusters
}

func (mc *MiniColumn) OutVector() bitarray.BitArray {
	return mc.outputVector
}

func (mc *MiniColumn) addNewClusters() {
	clusters := 0
	exist := 0
	activated := 0
	inputMax := 0
	outputMax := 0
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

			//inputMin := int(math.Sqrt(float64(3 * mc.inputLen)))
			//outputMin := int(math.Sqrt(float64(3 * mc.learningLen)))
			if receptorsActiveLen > inputMax {
				inputMax = receptorsActiveLen
			}
			if outputsActiveLen > outputMax {
				outputMax = outputsActiveLen
			}

			inputMin := 4
			outputMin := 4

			if receptorsActiveLen >= inputMin && outputsActiveLen >= outputMin && memorySize < mc.memoryLimit {
				cluster := NewCluster(activeReceptors, outputsActiveCount, mc.inputVectorLen, mc.outputVectorLen)
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
	fmt.Printf("Максимальная активность входа %d, выхода %d\n\n", inputMax, outputMax)
}
