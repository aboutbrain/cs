package cs

type HashMap map[string]bool

type CombinatorialSpace struct {
	LearningMode         int
	InternalTime         int
	NumberOfPoints       int
	NumberOfBitInPoint   uint64
	NumberOfBitInOutCode int
	Points               []Point
	OutHashSet           []HashMap
}

func NewCombinatorialSpace(size int, bit uint64, outCode int) *CombinatorialSpace {
	space := &CombinatorialSpace{NumberOfPoints: size, NumberOfBitInPoint: bit, NumberOfBitInOutCode: outCode}
	space.createPoints()
	space.OutHashSet = make([]HashMap, space.NumberOfBitInOutCode)
	for i := 0; i < space.NumberOfBitInOutCode; i++ {
		space.OutHashSet[i] = make(HashMap)
	}
	return space
}

func (cs *CombinatorialSpace) createPoints() {
	for i := 0; i < cs.NumberOfPoints; i++ {
		point := NewPoint(i, cs.NumberOfBitInPoint)
		point.OutBit = random(0, cs.NumberOfBitInOutCode)
		cs.Points = append(cs.Points, *point)
	}
}

func (cs *CombinatorialSpace) CheckOutHashSet(id int, hash string) bool {
	hashMap := cs.OutHashSet[id]
	if hashMap[hash] == true {
		return false
	}
	return true
}

func (cs *CombinatorialSpace) SetHash(id int, hash string) {
	cs.OutHashSet[id][hash] = true
}
