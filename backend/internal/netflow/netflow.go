package netflow

import (
	"time"
)

// NetFlowV5Header représente l'entête du format NetFlow v5
type NetFlowV5Header struct {
	Version      uint16
	Count        uint16
	SysUptime    uint32
	UnixSecs     uint32
	UnixNsecs    uint32
	FlowSequence uint32
	EngineType   uint8
	EngineID     uint8
	Sampling     uint16
}

// NetFlowV5Record représente un enregistrement unique de NetFlow v5
type NetFlowV5Record struct {
	SrcAddr   [4]byte
	DstAddr   [4]byte
	NextHop   [4]byte
	Input     uint16
	Output    uint16
	Packets   uint32
	Bytes     uint32
	StartTime uint32
	EndTime   uint32
	SrcPort   uint16
	DstPort   uint16
	Pad1      uint8
	TCPFlags  uint8
	Proto     uint8
	Tos       uint8
	SrcAS     uint16
	DstAS     uint16
	SrcMask   uint8
	DstMask   uint8
	Pad2      uint16
}

// ConvertNetFlowTime calcule le temps absolu d'un flux NetFlow
func ConvertNetFlowTime(header NetFlowV5Header, flowTime uint32) time.Time {
	baseTime := time.Unix(int64(header.UnixSecs), 0) // Temps en secondes basé sur UnixSecs
	flowStartTime := baseTime.Add(-time.Duration(header.SysUptime) * time.Millisecond).Add(time.Duration(flowTime) * time.Millisecond)
	return flowStartTime
}
