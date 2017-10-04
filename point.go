package cs

import "github.com/aboutbrain/cs/bitarray"

type Point struct {
	id               int
	inputVectorSize  int
	bitsPerInput     int
	receptorSet      bitarray.BitArray
	Memory           []Cluster
	potential        int
	outputVectorSize int
	OutputBit        int
}

func NewPoint(id, inputVectorSize int, bitsPerInput, outputVectorSize int) *Point {
	p := &Point{
		id:               id,
		inputVectorSize:  inputVectorSize,
		bitsPerInput:     int(bitsPerInput),
		outputVectorSize: outputVectorSize,
		receptorSet:      bitarray.NewBitArray(uint64(inputVectorSize)),
	}
	p.setReceptors()
	p.setOutBit()
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
		bitNumber := Random(0, p.inputVectorSize-1)
		if a, _ := p.receptorSet.GetBit(uint64(bitNumber)); a != true {
			p.receptorSet.SetBit(uint64(bitNumber))
		} else {
			goto rnd
		}
	}
}

func (p *Point) setOutBit() {
	outBitNum := Random(0, p.outputVectorSize-1)
	p.OutputBit = outBitNum
}

func (p *Point) GetReceptors() bitarray.BitArray {
	return p.receptorSet
}

func (p *Point) GetOutputBit() int {
	return p.OutputBit
}

func (p *Point) SetPotential(potential int) {
	p.potential = potential
}

func (p *Point) GetPotential() int {
	return p.potential
}
