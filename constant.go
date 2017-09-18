package cs

const (
	ModeSupervised = iota
	ModeUnsupervised
)

const (
	InputVectorSize            = 256
	OutputVectorSize           = 256
	ContextSize                = 10
	CombinatorialSpaceSize     = 60000
	BitsPerPoint               = 8
	ClusterThreshold           = 6
	ClusterActivationThreshold = 4
	CharacterBits              = 8
	PointMemoryCapacity        = 10
	PointContextCapacity       = 10
)
