package main

import (
	"fmt"

	"github.com/savreline/GoVector/govec/vclock"
)

func main() {
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
	// fmt.Println("C", clock.CopyFromMap(map1).Compare(map2, vclock.Concurrent))
	// fmt.Println("A", clock.CopyFromMap(map1).Compare(map2, vclock.Ancestor))
	// fmt.Println("D", clock.CopyFromMap(map1).Compare(map2, vclock.Descendant))
	// fmt.Println("E", clock.CopyFromMap(map1).Compare(map2, vclock.Equal))

	fmt.Println(clock.CopyFromMap(map1).CompareClocks(map2))
}

/*
	queueTest()
	printQueue()
	time.Sleep(3 * time.Second)
	ioutil.WriteFile("Repl"+noStr+".txt", []byte(eLog), 0644)

func queueTest() {
	map1 := map[string]uint64{
		"R1": 3,
		"R2": 2,
		"R3": 3,
	}
	opNode1 := OpNode{
		Type:      IK,
		Key:       "",
		Value:     "",
		Timestamp: map1,
		Pid:       noStr,
		ConcOp:    false}
	map2 := map[string]uint64{
		"R1": 2,
		"R2": 3,
		"R3": 3,
	}
	opNode2 := OpNode{
		Type:      IK,
		Key:       "",
		Value:     "",
		Timestamp: map2,
		Pid:       noStr,
		ConcOp:    false}
	map3 := map[string]uint64{
		"R1": 4,
		"R2": 4,
		"R3": 4,
	}
	opNode3 := OpNode{
		Type:      IK,
		Key:       "",
		Value:     "",
		Timestamp: map3,
		Pid:       noStr,
		ConcOp:    false}
	map4 := map[string]uint64{
		"R1": 2,
	}
	opNode4 := OpNode{
		Type:      IK,
		Key:       "",
		Value:     "",
		Timestamp: map4,
		Pid:       noStr,
		ConcOp:    false}
	map5 := map[string]uint64{
		"R2": 2,
	}
	opNode5 := OpNode{
		Type:      IK,
		Key:       "",
		Value:     "",
		Timestamp: map5,
		Pid:       noStr,
		ConcOp:    false}
	map6 := map[string]uint64{
		"R1": 6,
		"R2": 6,
		"R3": 6,
	}
	opNode6 := OpNode{
		Type:      IK,
		Key:       "",
		Value:     "",
		Timestamp: map6,
		Pid:       noStr,
		ConcOp:    false}
	map7 := map[string]uint64{
		"R1": 5,
		"R2": 5,
		"R3": 7,
	}
	opNode7 := OpNode{
		Type:      IK,
		Key:       "",
		Value:     "",
		Timestamp: map7,
		Pid:       noStr,
		ConcOp:    false}
	map8 := map[string]uint64{
		"R3": 2,
	}
	opNode8 := OpNode{
		Type:      IK,
		Key:       "",
		Value:     "",
		Timestamp: map8,
		Pid:       noStr,
		ConcOp:    false}
	addToQueue(opNode1)
	addToQueue(opNode2)
	addToQueue(opNode3)
	addToQueue(opNode4)
	addToQueue(opNode5)
	addToQueue(opNode6)
	addToQueue(opNode7)
	addToQueue(opNode8)
}
*/
