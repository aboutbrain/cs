package cs

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type Point struct {
	id             int
	bitsPerInput   int
	receptorSet    bitarray.BitArray
	bitsPerOutput  int
	outputSet      bitarray.BitArray
	Memory         []Cluster
	potential      int
	OutputBitArray []int
}

func NewPoint(id int, bitsPerInput, bitsPerOutput uint64) *Point {
	p := &Point{
		id:            id,
		bitsPerInput:  int(bitsPerInput),
		receptorSet:   bitarray.NewBitArray(bitsPerInput),
		bitsPerOutput: int(bitsPerOutput),
		outputSet:     bitarray.NewBitArray(bitsPerOutput),
	}
	p.setReceptors()
	p.setOutputs()
	return p
}

func (p *Point) SetMemory(cluster *Cluster) {
	p.Memory = append(p.Memory, *cluster)
}

func (p *Point) setReceptors() {
	for i := 0; i < p.bitsPerInput; i++ {
		bit := Random(0, p.bitsPerInput)
		p.receptorSet.SetBit(uint64(bit))
	}
}

func (p *Point) GetReceptors() bitarray.BitArray {
	return p.receptorSet
}

func (p *Point) setOutputs() {
	for i := 0; i < p.bitsPerInput; i++ {
		bit := Random(0, p.bitsPerInput)
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
