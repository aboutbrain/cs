package main

import (
	"fmt"

	"github.com/aboutbrain/cs"
    "github.com/golang-collections/go-datastructures/bitarray"

)

func main() {
	point := &cs.Point{}

	hipotesa := bitarray.NewBitArray(32)
	fmt.Println(hipotesa)
	hipotesa.SetBit(0)
	fmt.Printf("%d\n", hipotesa.ToNums()[0])
	hipotesa.SetBit(3)
	fmt.Printf("%d\n", hipotesa.ToNums()[1])
	hipotesa.SetBit(63)
	fmt.Printf("%d\n", hipotesa.ToNums()[2])
	point.Memory = append(point.Memory, hipotesa)

	hipotesa2 := bitarray.NewBitArray(32)
	fmt.Println(hipotesa)
	hipotesa2.SetBit(1)
	fmt.Printf("%d\n", hipotesa2.ToNums()[0])
	hipotesa2.SetBit(2)
	fmt.Printf("%d\n", hipotesa2.ToNums()[1])
	point.Memory = append(point.Memory, hipotesa2)

	iter := hipotesa.Blocks()
	iter.Next()
	a1, block1 := iter.Value()
	fmt.Println(a1, block1)
	iter.Next()
	a2, block2 := iter.Value()
	fmt.Println(a2, block2)
	fmt.Println(point)
}
