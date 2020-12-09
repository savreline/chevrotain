package main

import (
	"fmt"
	"testing"

	"../util"
	"github.com/savreline/GoVector/govec/vclock"
)

var maps []map[string]uint64
var opNodes []OpNode

func TestClockCmp(t *testing.T) {
	map1 := map[string]uint64{
		"R1": 2,
		"R2": 2,
		"R3": 3,
	}
	map2 := map[string]uint64{
		"R1": 2,
		"R2": 2,
		"R3": 2,
	}
	clock := vclock.New()
	fmt.Println("C", clock.CopyFromMap(map1).Compare(map2, vclock.Concurrent))
	fmt.Println("A", clock.CopyFromMap(map1).Compare(map2, vclock.Ancestor))
	fmt.Println("D", clock.CopyFromMap(map1).Compare(map2, vclock.Descendant))
	fmt.Println("E", clock.CopyFromMap(map1).Compare(map2, vclock.Equal))
	fmt.Println(clock.CopyFromMap(map1).CompareClocks(map2))
}

func Test123(t *testing.T) {
	initBasicExamples()
	for i := 0; i < 3; i++ {
		addToQueue(opNodes[i])
		printQueue()
		eLog = eLog + "\n"
	}
}

func Test321(t *testing.T) {
	initBasicExamples()
	for i := 2; i >= 0; i-- {
		addToQueue(opNodes[i])
		printQueue()
		eLog = eLog + "\n"
	}
}

func Test132(t *testing.T) {
	initBasicExamples()
	x := []int{0, 2, 1}
	for _, num := range x {
		addToQueue(opNodes[num])
		printQueue()
		eLog = eLog + "\n"
	}
}

func Test13524(t *testing.T) {
	initBasicExamples()
	x := []int{0, 2, 4, 1, 3}
	for _, num := range x {
		addToQueue(opNodes[num])
		printQueue()
		eLog = eLog + "\n"
	}
}

func Test24531(t *testing.T) {
	initBasicExamples()
	x := []int{1, 3, 4, 2, 0}
	for _, num := range x {
		addToQueue(opNodes[num])
		printQueue()
		eLog = eLog + "\n"
	}
}

func TestQueueAdv(t *testing.T) {
	initAdvExamples()
	for i := 0; i < 8; i++ {
		addToQueue(opNodes[i])
		printQueue()
		eLog = eLog + "\n"
	}
}

func TestConcQueue1(t *testing.T) {
	initConc1Examples()
	for i := 0; i < 8; i++ {
		addToQueue(opNodes[i])
		printQueue()
		eLog = eLog + "\n"
	}
}

func TestConcQueue2(t *testing.T) {
	initConc2Examples()
	for i := 0; i < 8; i++ {
		addToQueue(opNodes[i])
		printQueue()
		eLog = eLog + "\n"
	}
}

// initialization of examples
func initBasicExamples() {
	maps = make([]map[string]uint64, 5)
	opNodes = make([]OpNode, 5)
	for i := 0; i < 5; i++ {
		maps[i] = map[string]uint64{
			"R1": uint64(i + 1),
			"R2": 0,
			"R3": 0,
		}
		opNodes[i] = OpNode{
			OpType:    util.IK,
			Key:       "",
			Value:     "",
			Timestamp: maps[i],
			SrcPid:    noStr,
			ConcOp:    false}
	}
}

func initAdvExamples() {
	maps = make([]map[string]uint64, 8)
	opNodes = make([]OpNode, 8)
	maps[0] = map[string]uint64{
		"R1": 3,
		"R2": 2,
		"R3": 3,
	}
	maps[1] = map[string]uint64{
		"R1": 2,
		"R2": 3,
		"R3": 3,
	}
	maps[2] = map[string]uint64{
		"R1": 4,
		"R2": 4,
		"R3": 4,
	}
	maps[3] = map[string]uint64{
		"R1": 2,
	}
	maps[4] = map[string]uint64{
		"R2": 2,
	}
	maps[5] = map[string]uint64{
		"R1": 6,
		"R2": 6,
		"R3": 6,
	}
	maps[6] = map[string]uint64{
		"R1": 5,
		"R2": 5,
		"R3": 7,
	}
	maps[7] = map[string]uint64{
		"R3": 2,
	}
	for i := 0; i < 5; i++ {
		opNodes[i] = OpNode{
			OpType:    util.IK,
			Key:       "1",
			Value:     "1000",
			Timestamp: maps[i],
			SrcPid:    noStr,
			ConcOp:    false}
	}
	opNodes[5] = OpNode{
		OpType:    util.IV,
		Key:       "2",
		Value:     "2000",
		Timestamp: maps[5],
		SrcPid:    noStr,
		ConcOp:    false}
	opNodes[6] = OpNode{
		OpType:    util.RK,
		Key:       "3",
		Value:     "3000",
		Timestamp: maps[6],
		SrcPid:    noStr,
		ConcOp:    false}
	opNodes[7] = OpNode{
		OpType:    util.RV,
		Key:       "4",
		Value:     "5000",
		Timestamp: maps[7],
		SrcPid:    noStr,
		ConcOp:    false}
}

func initConc1Examples() {
	maps = make([]map[string]uint64, 8)
	opNodes = make([]OpNode, 8)
	maps[0] = map[string]uint64{
		"R1": 2,
		"R2": 2,
		"R3": 2,
	}
	maps[1] = map[string]uint64{
		"R1": 3,
		"R2": 2,
		"R3": 2,
	}
	maps[2] = map[string]uint64{
		"R1": 2,
		"R2": 3,
		"R3": 2,
	}
	maps[3] = map[string]uint64{
		"R1": 2,
		"R2": 2,
		"R3": 3,
	}
	maps[4] = map[string]uint64{
		"R1": 3,
		"R2": 3,
		"R3": 3,
	}
	maps[5] = map[string]uint64{
		"R1": 4,
		"R2": 3,
		"R3": 3,
	}
	maps[6] = map[string]uint64{
		"R1": 3,
		"R2": 4,
		"R3": 3,
	}
	maps[7] = map[string]uint64{
		"R1": 3,
		"R2": 3,
		"R3": 4,
	}
	for i := 0; i < 5; i++ {
		opNodes[i] = OpNode{
			OpType:    util.IK,
			Key:       "1",
			Value:     "1000",
			Timestamp: maps[i],
			SrcPid:    noStr,
			ConcOp:    false}
	}
	opNodes[5] = OpNode{
		OpType:    util.IV,
		Key:       "2",
		Value:     "2000",
		Timestamp: maps[5],
		SrcPid:    noStr,
		ConcOp:    false}
	opNodes[6] = OpNode{
		OpType:    util.RK,
		Key:       "3",
		Value:     "3000",
		Timestamp: maps[6],
		SrcPid:    noStr,
		ConcOp:    false}
	opNodes[7] = OpNode{
		OpType:    util.RV,
		Key:       "4",
		Value:     "5000",
		Timestamp: maps[7],
		SrcPid:    noStr,
		ConcOp:    false}
}

func initConc2Examples() {
	maps = make([]map[string]uint64, 8)
	opNodes = make([]OpNode, 8)
	maps[0] = map[string]uint64{
		"R1": 2,
		"R2": 2,
		"R3": 2,
	}
	maps[1] = map[string]uint64{
		"R1": 4,
		"R2": 1,
		"R3": 4,
	}
	maps[2] = map[string]uint64{
		"R1": 5,
		"R2": 1,
		"R3": 5,
	}
	maps[3] = map[string]uint64{
		"R1": 3,
		"R2": 3,
		"R3": 3,
	}
	maps[4] = map[string]uint64{
		"R1": 6,
		"R2": 6,
		"R3": 6,
	}
	maps[5] = map[string]uint64{
		"R1": 7,
		"R2": 7,
		"R3": 7,
	}
	maps[6] = map[string]uint64{
		"R1": 8,
		"R2": 8,
		"R3": 8,
	}
	maps[7] = map[string]uint64{
		"R1": 9,
		"R2": 9,
		"R3": 9,
	}
	for i := 0; i < 5; i++ {
		opNodes[i] = OpNode{
			OpType:    util.IK,
			Key:       "1",
			Value:     "1000",
			Timestamp: maps[i],
			SrcPid:    noStr,
			ConcOp:    false}
	}
	opNodes[5] = OpNode{
		OpType:    util.IV,
		Key:       "2",
		Value:     "2000",
		Timestamp: maps[5],
		SrcPid:    noStr,
		ConcOp:    false}
	opNodes[6] = OpNode{
		OpType:    util.RK,
		Key:       "3",
		Value:     "3000",
		Timestamp: maps[6],
		SrcPid:    noStr,
		ConcOp:    false}
	opNodes[7] = OpNode{
		OpType:    util.RV,
		Key:       "4",
		Value:     "5000",
		Timestamp: maps[7],
		SrcPid:    noStr,
		ConcOp:    false}
}
