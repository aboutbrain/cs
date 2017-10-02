package cs

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
	clustersPermanent1        int
	clustersPermanent2        int
}

func NewCombinatorialSpace(size, inputCodeSize int, receptors uint64, outCode int) *CombinatorialSpace {
	space := &CombinatorialSpace{
		NumberOfPoints:           size,
		NumberOfBitInInputCode: inputCodeSize,
		NumberOfReceptorsInPoint: receptors,
		NumberOfBitInOutCode:     outCode,
	}
	space.outBitToPointsMap = make(map[int][]int)
	space.createPoints()
	space.OutHashSet = make([]HashMap, space.NumberOfPoints)
	for i := 0; i < space.NumberOfPoints; i++ {
		space.OutHashSet[i] = make(HashMap)
	}
	return space
}

func (cs *CombinatorialSpace) createPoints() {
	for i := 0; i < cs.NumberOfPoints; i++ {
		point := NewPoint(i, cs.NumberOfBitInInputCode, int(cs.NumberOfReceptorsInPoint), cs.NumberOfBitInOutCode)
		cs.Points = append(cs.Points, *point)
		outBit := point.GetOutputBit()
		arr := cs.outBitToPointsMap[outBit]
		arr = append(arr, i)
		cs.outBitToPointsMap[outBit] = arr
	}
}

func (cs *CombinatorialSpace) DeleteCluster(point *Point, clusterId int) {
	cluster := point.Cluster(clusterId)
	status := cluster.Status
	point.DeleteCluster(clusterId)
	cs.clustersTotal--
	if status == ClusterPermanent1 {
		cs.clustersPermanent1--
	}
	if status == ClusterPermanent2 {
		cs.clustersPermanent2--
	}
	cs.RemoveHash(point.id, cluster.GetHash())
}

func (cs *CombinatorialSpace) GetPointsByOutBitNumber(id int) []int {
	return cs.outBitToPointsMap[id]
}

func (cs *CombinatorialSpace) CheckOutHashSet(pointId int, hash string) bool {
	hashMap := cs.OutHashSet[pointId]
	_, ok := hashMap[hash]
	return ok
}

func (cs *CombinatorialSpace) SetHash(pointId int, hash string) {
	cs.OutHashSet[pointId][hash] = true
}

func (cs *CombinatorialSpace) RemoveHash(pointId int, hash string) {
	delete(cs.OutHashSet[pointId], hash)
}

func (cs *CombinatorialSpace) IncreaseClusters() {
	cs.clustersTotal++
}

func (cs *CombinatorialSpace) ClustersCounters() (int, int, int) {
	return cs.clustersTotal, cs.clustersPermanent1, cs.clustersPermanent2
}
