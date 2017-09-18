package cs

import "github.com/golang-collections/go-datastructures/bitarray"

type Context struct {
	ConceptId int
}

type Concept struct {
	ContextId int
}

type Point struct {
	Memory []bitarray.BitArray
}

func (p *Point) SetMemory(hypothesis *bitarray.BitArray) {
	if true {
		p.Memory = append(p.Memory, *hypothesis)
	}
}
