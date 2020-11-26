package main

import (
	"fmt"

	"github.com/savreline/GoVector/govec/vclock"
)

var maps []map[string]uint64
var opNodes []OpNode

/* To be pasted into main
// queueTest123()
// queueTest321()
// queueTest132()
// queueTest13524()
// queueTest24531()
queueTestAdv()
processConcOps()
printQueue()
fmt.Println(eLog)
os.Exit(0)
*/

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
			Type:      IK,
			Key:       "",
			Value:     "",
			Timestamp: maps[i],
			Pid:       noStr,
			ConcOp:    false}
	}
}

func queueTest123() {
	initBasicExamples()
	for i := 0; i < 3; i++ {
		addToQueue(opNodes[i])
		printQueue()
		eLog = eLog + "\n"
	}
}

func queueTest321() {
	initBasicExamples()
	for i := 2; i >= 0; i-- {
		addToQueue(opNodes[i])
		printQueue()
		eLog = eLog + "\n"
	}
}

func queueTest132() {
	initBasicExamples()
	x := []int{1, 3, 2}
	for i := range x {
		addToQueue(opNodes[i])
		printQueue()
		eLog = eLog + "\n"
	}
}

func queueTest13524() {
	initBasicExamples()
	x := []int{1, 3, 5, 2, 4}
	for i := range x {
		addToQueue(opNodes[i])
		printQueue()
		eLog = eLog + "\n"
	}
}

func queueTest24531() {
	initBasicExamples()
	x := []int{2, 4, 5, 3, 1}
	for i := range x {
		addToQueue(opNodes[i])
		printQueue()
		eLog = eLog + "\n"
	}
}

func clockCmpTest() {
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

// Advanced Test
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
			Type:      IK,
			Key:       "1",
			Value:     "1000",
			Timestamp: maps[i],
			Pid:       noStr,
			ConcOp:    false}
	}
	opNodes[5] = OpNode{
		Type:      IV,
		Key:       "2",
		Value:     "2000",
		Timestamp: maps[5],
		Pid:       noStr,
		ConcOp:    false}
	opNodes[6] = OpNode{
		Type:      RK,
		Key:       "3",
		Value:     "3000",
		Timestamp: maps[6],
		Pid:       noStr,
		ConcOp:    false}
	opNodes[7] = OpNode{
		Type:      RV,
		Key:       "4",
		Value:     "5000",
		Timestamp: maps[7],
		Pid:       noStr,
		ConcOp:    false}
}

func queueTestAdv() {
	initAdvExamples()
	for i := 0; i < 8; i++ {
		addToQueue(opNodes[i])
		printQueue()
		eLog = eLog + "\n"
	}
}
