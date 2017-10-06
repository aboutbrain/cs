package cs

import "log"

type HashMap map[string]bool

type CombinatorialSpace struct {
	LearningMode             int
	InternalTime             int
	NumberOfPoints           int
	NumberOfBitInInputCode   int
	NumberOfReceptorsInPoint uint64
	NumberOfBitInOutCode     int
	Points                   []Point
	OutHashSet               []HashMap
	outBitToPointsMap        map[int][]int
	clustersTotal            int
	clustersPermanent1       int
	clustersPermanent2       int
}

func NewCombinatorialSpace(size, inputCodeSize int, receptors uint64, outCode int) *CombinatorialSpace {
	space := &CombinatorialSpace{
		NumberOfPoints:           size,
		NumberOfBitInInputCode:   inputCodeSize,
		NumberOfReceptorsInPoint: receptors,
		NumberOfBitInOutCode:     outCode,
		Points:                   make([]Point, size, size),
	}
	space.outBitToPointsMap = make(map[int][]int)
	space.createPoints()
	space.OutHashSet = make([]HashMap, space.NumberOfBitInOutCode)
	for i := 0; i < space.NumberOfBitInOutCode; i++ {
		space.OutHashSet[i] = make(HashMap)
	}
	return space
}

func (cs *CombinatorialSpace) createPoints() {
	for i := 0; i < cs.NumberOfPoints; i++ {
		point := NewPoint(i, cs.NumberOfBitInInputCode, int(cs.NumberOfReceptorsInPoint), cs.NumberOfBitInOutCode)
		//cs.Points = append(cs.Points, *point)
		cs.Points[i] = *point
		outBit := point.GetOutputBit()
		arr := cs.outBitToPointsMap[outBit]
		arr = append(arr, i)
		cs.outBitToPointsMap[outBit] = arr
	}
}

func (cs *CombinatorialSpace) DeleteCluster(point *Point, clusterId int, hashRemove bool) {
	cluster := point.Cluster(clusterId)
	status := cluster.Status
	hash := cluster.GetHash()
	point.DeleteCluster(clusterId)
	cs.clustersTotal--
	if status == ClusterPermanent1 {
		cs.clustersPermanent1--
	}
	if status == ClusterPermanent2 {
		cs.clustersPermanent2--
	}
	if hashRemove {
		cs.RemoveHash(point.OutputBit, hash)
	}
}

func (cs *CombinatorialSpace) GetPointsByOutBitNumber(id int) []int {
	return cs.outBitToPointsMap[id]
}

func (cs *CombinatorialSpace) CheckOutHashSet(outBit int, hash string) bool {
	hashMap := cs.OutHashSet[outBit]
	_, ok := hashMap[hash]
	return ok
}

func (cs *CombinatorialSpace) SetHash(outBit int, hash string) {
	cs.OutHashSet[outBit][hash] = true
}

func (cs *CombinatorialSpace) RemoveHash(outBit int, hash string) {
	val, ok := cs.OutHashSet[outBit][hash]
	if !ok {
		log.Printf("нет такого элемента (%s) в карте, value: %b", hash, val)
	}
	delete(cs.OutHashSet[outBit], hash)
}

func (cs *CombinatorialSpace) IncreaseClusters() {
	cs.clustersTotal++
}

func (cs *CombinatorialSpace) ClustersCounters() (int, int, int) {
	return cs.clustersTotal, cs.clustersPermanent1, cs.clustersPermanent2
}
