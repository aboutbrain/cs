package cs

import "github.com/golang-collections/go-datastructures/bitarray"

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

func (p *Point) setReceptors() {
	for i := 0; i < p.bitsPerInput; i++ {
		bit := Random(0, p.inputVectorSize)
		p.receptorSet.SetBit(uint64(bit))
	}
}

func (p *Point) GetReceptors() bitarray.BitArray {
	return p.receptorSet
}

func (p *Point) setOutputs() {
	for i := 0; i < p.bitsPerOutput; i++ {
		bit := Random(0, p.outputVectorSize)
		p.outputSet.SetBit(uint64(bit))
		p.OutputBitArray = append(p.OutputBitArray, bit)
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
