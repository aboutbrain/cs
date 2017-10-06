package persist

import (
	"encoding/json"
	"io/ioutil"

	"github.com/aboutbrain/cs/bitarray"
	"github.com/aboutbrain/cs/text"
)

type MapOfBits map[int][]uint64
type CodesDump struct {
	Capacity  uint64
	CharCodes map[int]MapOfBits
}

func ToFile(path string, n *text.CharContextCodes) {
	dump := ToDump(n)
	DumpToFile(path, dump)
}

func DumpToFile(path string, dump *CodesDump) {
	j, _ := json.Marshal(dump)
	err := ioutil.WriteFile(path, j, 0644)
	if err != nil {
		panic(err)
	}
}

func ToDump(n *text.CharContextCodes) *CodesDump {
	dump := &CodesDump{Capacity: uint64(n.VectorCapacity), CharCodes: make(map[int]MapOfBits)}
	for i, v := range n.CharContext {
		dump.CharCodes[i] = MapOfBits{}
		for j, v1 := range v {
			nums := v1.ToNums()
			dump.CharCodes[i][j] = nums
		}
	}
	return dump
}

func DumpFromFile(path string) *CodesDump {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	dump := &CodesDump{}
	err = json.Unmarshal(b, dump)
	if err != nil {
		panic(err)
	}
	return dump
}

func FromFile(path string) *text.CharContextCodes {
	dump := DumpFromFile(path)
	n := FromDump(dump)
	return n
}

func FromDump(dump *CodesDump) *text.CharContextCodes {
	codes := make(text.CharContext)
	for i, v := range dump.CharCodes {
		arr := make(map[int]bitarray.BitArray)
		for j, v1 := range v {
			contArray := bitarray.NewBitArray(dump.Capacity)
			for _, v2 := range v1 {
				contArray.SetBit(v2)
			}
			arr[j] = contArray
		}
		codes[i] = arr
	}
	charContextCode := text.CharContextCodes{VectorCapacity: int(dump.Capacity), CharContext: codes}
	return &charContextCode
}
