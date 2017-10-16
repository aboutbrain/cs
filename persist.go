package cs

import (
	"encoding/json"
	"io/ioutil"
)

type ClusterDump struct {
	Status        int
	StartTime     int
	InputWeights  map[int]float32
	OutputWeights map[int]float32
	NLearn        float64
	SXY           float32
	SX            float32
	SY            float32
	RealSXY       float32
}

type PointDump struct {
	Id          int
	ReceptorSet []uint64
	OutputSet   []uint64
	Memory      []ClusterDump
}

type CombinatorialSpaceDump struct {
	InternalTime int
	PointsDump   []PointDump
}

type MiniColumnDump struct {
	InputVectorLen             int
	OutputVectorLen            int
	ClusterThreshold           int
	ClusterActivationThreshold int
	PointMemoryLimit           int
	OutBitActivationLevel      int
	CombinatorialSpace         CombinatorialSpaceDump
}

func ToFile(path string, column *MiniColumn) {
	dump := ToDump(column)
	DumpToFile(path, dump)
}

func DumpToFile(path string, dump *MiniColumnDump) {
	j, _ := json.MarshalIndent(dump, "", "\t")
	err := ioutil.WriteFile(path, j, 0644)
	if err != nil {
		panic(err)
	}
}

func ToDump(mc *MiniColumn) *MiniColumnDump {
	dump := &MiniColumnDump{
		InputVectorLen:             int(mc.inputVectorLen),
		OutputVectorLen:            int(mc.outputVectorLen),
		ClusterThreshold:           mc.clusterThreshold,
		ClusterActivationThreshold: mc.clusterActivationThreshold,
		PointMemoryLimit:           mc.memoryLimit,
		OutBitActivationLevel:      mc.level,
	}

	pointsDump := make([]PointDump, len(mc.cs.Points))

	for pointId, point := range mc.cs.Points {
		pDump := PointDump{
			Id:          point.id,
			ReceptorSet: point.receptorSet.ToNums(),
			OutputSet:   point.outputSet.ToNums(),
		}
		pMemory := make([]ClusterDump, len(point.Memory))
		for clusterId, cluster := range point.Memory {
			cDump := ClusterDump{
				Status:        cluster.Status,
				StartTime:     cluster.startTime,
				InputWeights:  cluster.inputWeights,
				OutputWeights: cluster.outputWeights,
				NLearn:        cluster.nLearn,
				SXY:           cluster.sXY,
				SX:            cluster.sX,
				SY:            cluster.sY,
				RealSXY:       cluster.realSXY,
			}
			pMemory[clusterId] = cDump
		}
		pDump.Memory = pMemory
		pointsDump[pointId] = pDump
	}

	csDump := CombinatorialSpaceDump{
		InternalTime: mc.cs.InternalTime,
		PointsDump:   pointsDump,
	}

	dump.CombinatorialSpace = csDump

	return dump
}
