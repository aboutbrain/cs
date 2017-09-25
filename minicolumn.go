package cs

import "github.com/golang-collections/go-datastructures/bitarray"

type MiniColumn struct {
	inputVector      bitarray.BitArray
	outputVector     bitarray.BitArray
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
		level:            2,
	}
}

func (mc *MiniColumn) SetCombinatorialSpace(cs *CombinatorialSpace) {
	mc.cs = cs
}

func (mc *MiniColumn) SetInputVector(inputVector bitarray.BitArray) {
	mc.inputVector = inputVector
}

func (mc *MiniColumn) SetLearningVector(learningVector bitarray.BitArray) {
	mc.learningVector = learningVector
}

func (mc *MiniColumn) Next() {
	mc.ActivateClusters()
	mc.MakeOutVector()
	/*mc.ModifyClusters()
	mc.ConsolidateMemory()*/
}

func (mc *MiniColumn) ActivateClusters() {
	for i, point := range mc.cs.Points {
		for j, cluster := range point.Memory {
			cluster.SetCurrentPotential(1)
			point.Memory[j] = cluster
		}
		point.SetPotential(2)
		mc.cs.Points[i] = point
	}
}

func (mc *MiniColumn) MakeOutVector() {
	for i, currentOutBitPointsMap := range mc.cs.outBitToPointsMap {
		p := 0
		for _, pointId := range currentOutBitPointsMap {
			p += mc.cs.Points[pointId].GetPotential()
		}
		if p > mc.level {
			mc.outputVector.SetBit(uint64(i))
		}
	}
}

func (mc *MiniColumn) AddNewClusters() {
	for _, v := range mc.learningVector.ToNums() {
		points := mc.cs.GetPointsByOutBitNumber(int(v))
		for _, pointId := range points {
			p := mc.cs.Points[pointId]
			cluster := NewCluster(mc.inputVector, p.GetReceptors())
			hash := cluster.GetHash()
			size := cluster.GetSize()
			memorySize := len(p.Memory)
			if mc.cs.CheckOutHashSet(p.OutBit, hash) && size >= mc.clusterThreshold && memorySize < mc.memoryLimit {
				mc.cs.SetHash(p.OutBit, hash)
				p.SetMemory(cluster)
				mc.cs.Points[pointId] = p
			}
		}
	}
	mc.epoch++
}
