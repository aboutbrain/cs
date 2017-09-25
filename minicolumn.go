package cs

import "github.com/golang-collections/go-datastructures/bitarray"

type MiniColumn struct {
	inputVector      bitarray.BitArray
	outputVector     bitarray.BitArray
	learningVector   bitarray.BitArray
	cs               *CombinatorialSpace
	clusterThreshold int
	epoch            int
}

func NewMiniColumn(clusterThreshold int) *MiniColumn {
	return &MiniColumn{clusterThreshold: clusterThreshold}
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

func (mc *MiniColumn) AddNewClusters() {
	for _, v := range mc.learningVector.ToNums() {
		points := mc.cs.GetPointsByOutBitNumber(int(v))
		for _, pointId := range points {
			p := mc.cs.Points[pointId]
			cluster := NewCluster(mc.inputVector, p.GetReceptors())
			hash := cluster.GetHash()
			size := cluster.GetSize()
			if mc.cs.CheckOutHashSet(p.OutBit, hash) && size >= mc.clusterThreshold {
				mc.cs.SetHash(p.OutBit, hash)
				p.SetMemory(cluster)
				mc.cs.Points[pointId] = p
			}
		}
	}
	mc.epoch++
}
