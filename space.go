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

func (cs *CombinatorialSpace) DeleteCluster(pointId, clusterId int) {
	point := cs.Points[pointId]
	cluster := point.Cluster(clusterId)
	point.DeleteCluster(clusterId)
	cs.clustersTotal--
	if cluster.Status == ClusterPermanent2 {
		cs.clustersPermanent--
	}
}

func (cs *CombinatorialSpace) GetPointsByOutBitNumber(id int) []int {
	return cs.outBitToPointsMap[id]
}

func (cs *CombinatorialSpace) CheckOutHashSet(pointId int, hash string) bool {
	hashMap := cs.OutHashSet[pointId]
	if hashMap[hash] == true {
		return false
	}
	return true
}

func (cs *CombinatorialSpace) SetHash(id int, hash string) {
	cs.OutHashSet[id][hash] = true
	cs.clustersTotal++
}

func (cs *CombinatorialSpace) IncreaseClusters() {
	cs.clustersTotal++
}

func (cs *CombinatorialSpace) ClustersCounters() (int, int) {
	return cs.clustersTotal, cs.clustersPermanent
}
