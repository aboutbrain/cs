package main

import (
	"fmt"

	"github.com/aboutbrain/cs"
)

func main() {
	point := &cs.Point{}

	array := [8]byte{0,0,0,0,0,0,0,0}
	fmt.Println(array)

	point.Memory = append(point.Memory, array)
	array = [8]byte{0,0,0,0,0,0,0,1}
	point.Memory = append(point.Memory, array)
	array = [8]byte{0,0,0,0,0,0,1,0}
	point.Memory = append(point.Memory, array)

	/*point.InputArrayMap[1] = 1
	point.InputArrayMap[2] = 100
	conceptId := 1
	contextId := 2

		concept := cs.Concept{contextId}
		context := cs.Context{conceptId}	/*point.Concept[conceptId] = concept
		point.Context[contextId] = context*/
	fmt.Println(point)
}
