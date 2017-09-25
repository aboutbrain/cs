package cs

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type Point struct {
	id           int
	bitsPerPoint int
	receptorSet  bitarray.BitArray
	Memory       []Cluster
	potential    int
	OutBit       int
}

func NewPoint(id int, bitsPerPoint uint64) *Point {
	p := &Point{
		id:           id,
		bitsPerPoint: int(bitsPerPoint),
		receptorSet:  bitarray.NewBitArray(bitsPerPoint),
	}
	p.setReceptors()
	return p
}

func (p *Point) SetMemory(cluster *Cluster) {
	p.Memory = append(p.Memory, *cluster)
}

func (p *Point) setReceptors() {
	for i := 0; i < p.bitsPerPoint; i++ {
		bit := random(0, p.bitsPerPoint)
		p.receptorSet.SetBit(uint64(bit))
	}
}

func (p *Point) GetReceptors() bitarray.BitArray {
	return p.receptorSet
}

func (p *Point) SetPotential(i int) {
	p.potential = i
}

func (p *Point) GetPotential() int{
	return p.potential
}
