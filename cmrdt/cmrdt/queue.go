package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"../../util"
	"github.com/savreline/GoVector/govec/vclock"
)

// OpCode is an operation code
type OpCode int

// OpCodes
const (
	IK OpCode = iota + 1
	IV
	RK
	RV
)

// OpNode represents a node in the operation wait queue
type OpNode struct {
	Type       OpCode
	Key, Value string
	Timestamp  vclock.VClock
	Pid        string
	ConcOp     bool
}

// translate operation code from string to op code
func lookupOpCode(opName string) OpCode {
	if opName == "IK" {
		return IK
	} else if opName == "IV" {
		return IV
	} else if opName == "RK" {
		return RK
	} else if opName == "RV" {
		return RV
	} else {
		util.PrintErr(noStr, errors.New("lookupOpCode: unknown operation"))
		return 0
	}
}

// Print the Queue
func printQueue() {
	lock.Lock()
	for n := queue.Front(); n != nil; n = n.Next() {
		eLog = eLog + fmt.Sprintln(n.Value)
	}
	eLog = eLog + "\n"
	lock.Unlock()
}

// insert a node into the correct location in the queue
func addToQueue(node OpNode) {
	lock.Lock()
	if queue.Front() == nil {
		queue.PushFront(node)
		lock.Unlock()
		return
	}
	for curNode := queue.Front(); curNode != nil; curNode = curNode.Next() {
		cmp := node.Timestamp.CompareClocks(curNode.Value.(OpNode).Timestamp)

		if cmp == 1 {
			node.ConcOp = true
			queue.InsertBefore(node, curNode)
			lock.Unlock()
			return
		}
		if cmp == 2 {
			queue.InsertBefore(node, curNode)
			lock.Unlock()
			return
		}
	}
	queue.PushBack(node)
	lock.Unlock()
}

// process some of the operations that are queued up
func processQueue() {
	for {
		time.Sleep(2 * time.Second)
		lock.Lock()
		eliminateConcOps()
		processQueueHelper()
		lock.Unlock()
	}
}

// processQueueHelper does the actual processing of queue operations
func processQueueHelper() {
	updateCurTick()

	for n := queue.Front(); n != nil; {
		opNode := n.Value.(OpNode)

		/* Stop if any timestamp is exceeding the current safe tick */
		for i := 0; i < len(conns); i++ {
			if int(opNode.Timestamp["R"+strconv.Itoa(i+1)]) > curTick {
				return
			}
		}

		/* Run the associated op */
		if opNode.Type == IK {
			InsertKeyLocal(opNode.Key)
		} else if opNode.Type == IV {
			InsertValueLocal(opNode.Key, opNode.Value)
		}

		/* Delete the associated node */
		temp := n
		n = n.Next()
		queue.Remove(temp)
	}
}

// updateCurTick updates the current "safe" tick
func updateCurTick() {
	noReplicas := len(conns)

	/* Gather all current timestamps */
	a := make([][]int, noReplicas)
	for n := queue.Front(); n != nil; n = n.Next() {
		for i := 0; i < noReplicas; i++ {
			a[i] = append(a[i], int(n.Value.(OpNode).Timestamp["R"+strconv.Itoa(i+1)]))
		}
	}

	if len(a[0]) == 0 {
		return
	}

	/* Determine the latest timestamp per replica */
	b := make([]int, noReplicas)
	for i := 0; i < noReplicas; i++ {
		b[i] = max(a, i)
	}

	/* Determine the earliest timestamp for all replicas */
	curTick = min(b)
	eLog = eLog + fmt.Sprint(a) + ":" + fmt.Sprint(curTick) + "\n"
}

// determine the latest timestamp per replica
func max(a [][]int, i int) int {
	sort.Ints(a[i])
	for j := 0; j < len(a[i])-1; j++ {
		if a[i][j+1] > curTick+1 && a[i][j+1] > a[i][j]+1 {
			return a[i][j]
		}
	}
	return a[i][len(a[i])-1]
}

// determine the earliest timestamp for all replicas
func min(b []int) int {
	res := b[0]
	for j := 1; j < len(b); j++ {
		if b[j] < res {
			res = b[j]
		}
	}
	return res
}

// eliminate concurrent operations from the queue using predefined preference
func eliminateConcOps() {
	// TODO
}
