package cs

type CombinatorialSpace struct {
	LearningMode         int
	InternalTime         int
	NumberOfPoints       int
	NumberOfBitInPoint   int
	NumberOfBitInOutCode int
	Points                [CombinatorialSpaceSize]Point
}
