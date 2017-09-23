package cs

type CombinatorialSpace struct {
	LearningMode         int
	InternalTime         int
	NumberOfPoints       int
	NumberOfBitInPoint   uint64
	NumberOfBitInOutCode int
	Points               []Point
}

func NewCombinatorialSpace(size int, bit uint64) *CombinatorialSpace {
	space := &CombinatorialSpace{NumberOfPoints: size, NumberOfBitInPoint: bit}
	space.createPoints()
	return space
}

func (cs *CombinatorialSpace) createPoints() {
	for i := 0; i < cs.NumberOfPoints; i++ {
		point := NewPoint(i, cs.NumberOfBitInPoint)
		//point.setReceptors()
		cs.Points = append(cs.Points, *point)
	}
}
