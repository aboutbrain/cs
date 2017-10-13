package cs

import "github.com/aboutbrain/cs/bitarray"

type Point struct {
	id               int
	inputVectorSize  int
	outputVectorSize int
	bitsPerInput     int
	receptorSet      bitarray.BitArray
	bitsPerOutput    int
	outputSet        bitarray.BitArray
	Memory           []Cluster
	potential        int
	activated        int
	OutputBitArray   []int
}

func NewPoint(id, inputVectorSize, outputVectorSize int, bitsPerInput, bitsPerOutput uint64) *Point {
	p := &Point{
		id:               id,
		inputVectorSize:  inputVectorSize,
		outputVectorSize: outputVectorSize,
		bitsPerInput:     int(bitsPerInput),
		receptorSet:      bitarray.NewBitArray(uint64(inputVectorSize)),
		bitsPerOutput:    int(bitsPerOutput),
		outputSet:        bitarray.NewBitArray(uint64(outputVectorSize)),
	}
	p.setReceptors()
	p.setOutputs()
	return p
}

func (p *Point) GetId() int {
	return p.id
}

func (p *Point) SetMemory(cluster *Cluster) {
	p.Memory = append(p.Memory, *cluster)
}

func (p *Point) Cluster(clusterId int) *Cluster {
	return &p.Memory[clusterId]
}

func (p *Point) DeleteCluster(clusterId int) *Point {
	p.Memory = append(p.Memory[:clusterId], p.Memory[clusterId+1:]...)
	return p
}

func (p *Point) setReceptors() {
	for i := 0; i < p.bitsPerInput; i++ {
	rnd:
		bitNumber := Random(0, p.inputVectorSize)
		if a, _ := p.receptorSet.GetBit(uint64(bitNumber)); a != true {
			p.receptorSet.SetBit(uint64(bitNumber))
		} else {
			goto rnd
		}
	}
}

func (p *Point) GetReceptors() bitarray.BitArray {
	return p.receptorSet
}

func (p *Point) setOutputs() {
	for i := 0; i < p.bitsPerOutput; i++ {
	rnd:
		bitNumber := Random(0, p.outputVectorSize)
		if a, _ := p.outputSet.GetBit(uint64(bitNumber)); a != true {
			p.outputSet.SetBit(uint64(bitNumber))
			p.OutputBitArray = append(p.OutputBitArray, bitNumber)
		} else {
			goto rnd
		}
	}
}

func (p *Point) GetOutputs() bitarray.BitArray {
	return p.outputSet
}

func (p *Point) SetPotential(i int) {
	p.potential = i
}

func (p *Point) GetPotential() int {
	return p.potential
}
