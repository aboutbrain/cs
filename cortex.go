package cs

import (
	"fmt"

	"github.com/aboutbrain/cs/bitarray"
)

type CortexContext struct {
	Id           int
	InputVector  bitarray.BitArray
	OutputVector bitarray.BitArray
	TextFragment string
	p            int
}

type Cortex struct {
	mc           *MiniColumn
	capacity     int
	ContextData  []CortexContext
	codes        *CharContextCodes
	textFragment string
	wordID       int
}

func NewCortex(mc *MiniColumn, capacity int, codes *CharContextCodes) *Cortex {
	return &Cortex{
		mc:          mc,
		capacity:    capacity,
		ContextData: make([]CortexContext, capacity),
		codes:       codes,
	}
}

func (c *Cortex) Run(contextData *[]CortexContext) {
	c.ContextData = *contextData
	c.MakeOutCode()
	fmt.Printf("Ищем победителя... \n")
	winnerId := c.Winner()
	fmt.Printf("Победитель контекст: %d\n", winnerId)
	c.LearnUnsupervised(winnerId)
	c.QualityControl()
}

func (c *Cortex) MakeOutCode() {
	for i := range c.ContextData {
		c.mc.SetInputVector(c.ContextData[i].InputVector)
		c.mc.InputText = c.ContextData[i].TextFragment
		fmt.Printf("i: %d, ContextID: %d, InputText  : \"%s\"\n", c.wordID, i, c.ContextData[i].TextFragment)
		c.ContextData[i].OutputVector, c.ContextData[i].p = c.mc.CreateOutput()
		c.showVectors(c.ContextData[i].InputVector, c.ContextData[i].OutputVector)
		total := c.mc.cs.ClustersCounters()
		fmt.Printf("Clusters: %d, сработавших: %d\n\n", total, c.ContextData[i].p)
	}
}

func (c *Cortex) Winner() int {
	max := 0
	iMax := 0
	for i, j := range c.ContextData {
		if j.p > max {
			max = j.p
			iMax = i
		}
	}
	return iMax
}

func (c *Cortex) LearnUnsupervised(id int) {
	fmt.Printf("Обучаем контекст: %d, текст: \"%s\"\n", id, c.ContextData[id].TextFragment)
	c.mc.LearnUnsupervised(c.ContextData[id].InputVector, c.ContextData[id].OutputVector)
}

func (c *Cortex) QualityControl() {
	fmt.Printf("Контроль качества кода!\n\n")
}

func (c *Cortex) showVectors(source, output bitarray.BitArray) {
	input, inputNum := BitArrayToString(source, int(c.mc.inputVectorLen))
	fmt.Printf("InputVector:   %s, %d\n", input, inputNum)
	outputStr, outputNum := BitArrayToString(output, int(c.mc.outputVectorLen))
	fmt.Printf("OutputVector:  %s, %d\n", outputStr, outputNum)

	if outputNum > 8 {
		fmt.Println("\033[31mFAIL!!\033[0m")
	}
	//fmt.Printf("LerningVector: %s\n", BitArrayToString(learning, int(c.mc.outputVectorLen)))
	//fmt.Printf("DeltaVector:   %s\n", BitArrayToString2(output, learning, int(c.mc.outputVectorLen)))
	/*
		if len(output.ToNums()) > 0 {
			if !nVector {
				fmt.Println("\033[31mFAIL!!\033[0m")
			} else {
				fmt.Println("\033[32mPASS!!\033[0m")
			}
		}*/
}

func BitArrayToString2(output, learning bitarray.BitArray, vectorLen int) string {
	delta := output.And(learning)
	nums := delta.ToNums()
	s := ""
	for i := 0; i < vectorLen; i++ {
		if InArray64(vectorLen-1-i, nums) {
			s += "\033[32m1\033[0m"
		} else {
			s += "0"
		}
	}
	return s
}
