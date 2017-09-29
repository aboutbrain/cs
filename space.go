package cs

type HashMap map[string]bool

type CombinatorialSpace struct {
	LearningMode             int
	InternalTime             int
	NumberOfPoints           int
	NumberOfReceptorsInPoint uint64
	NumberOfOutputsInPoint   uint64
	NumberOfBitInOutCode     int
	Points                   []Point
	OutHashSet               []HashMap
	outBitToPointsMap        map[int][]int
	clustersTotal            int
	clustersPermanent        int
}

func NewCombinatorialSpace(size int, receptors, outputs uint64, outCode int) *CombinatorialSpace {
	space := &CombinatorialSpace{
		NumberOfPoints:           size,
		NumberOfReceptorsInPoint: receptors,
		NumberOfOutputsInPoint:   outputs,
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
		point := NewPoint(i, cs.NumberOfBitInOutCode, cs.NumberOfBitInOutCode, cs.NumberOfReceptorsInPoint, cs.NumberOfOutputsInPoint)
		cs.Points = append(cs.Points, *point)
		outBits := point.GetOutputs()
		for _, v := range outBits.ToNums() {
			arr := cs.outBitToPointsMap[int(v)]
			arr = append(arr, i)
			cs.outBitToPointsMap[int(v)] = arr
		}
	}
}

func (cs *CombinatorialSpace) DeleteCluster(point *Point, clusterId int) {
	cluster := point.Cluster(clusterId)
	status := cluster.Status
	point.DeleteCluster(clusterId)
	cs.clustersTotal--
	if status == ClusterPermanent2 {
		cs.clustersPermanent--
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

func (cs *CombinatorialSpace) ClustersCounters() (int, int) {
	return cs.clustersTotal, cs.clustersPermanent
}
