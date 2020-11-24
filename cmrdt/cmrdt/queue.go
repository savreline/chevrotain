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

// ListNode is a node in the linked list queue
type ListNode struct {
	Data OpNode
	Next *ListNode
}

// translate operation code from string to op code
func lookupOpCode(opCode OpCode) string {
	if opCode == IK {
		return "IK"
	} else if opCode == IV {
		return "IV"
	} else if opCode == RK {
		return "RK"
	} else if opCode == RV {
		return "RV"
	} else {
		util.PrintErr(noStr, errors.New("lookupOpCode: unknown operation"))
		return ""
	}
}

// Print the Queue
func printQueue() {
	for n := queue; n != nil; n = n.Next {
		eLog = eLog + fmt.Sprintln(n.Data)
	}
}

// insert a node into the correct location in the queue
func addToQueue(node OpNode) {
	lock.Lock()

	/* Case 1: Empty Queue */
	if queue == nil {
		queue = &ListNode{Data: node, Next: nil}
		lock.Unlock()
		return
	}

	/* Case 2: Insertion at the Head */
	cmp := node.Timestamp.CompareClocks(queue.Data.Timestamp)
	if cmp == 1 {
		node.ConcOp = true
	}
	if cmp == 1 || cmp == 2 {
		queue = &ListNode{Data: node, Next: queue}
		lock.Unlock()
		return
	}

	/* Case 3: Insertion in the Middle */
	curNode := queue
	for ; curNode.Next != nil; curNode = curNode.Next {
		cmp := node.Timestamp.CompareClocks(curNode.Next.Data.Timestamp)

		if cmp == 1 {
			node.ConcOp = true
		}
		if cmp == 1 || cmp == 2 {
			curNode.Next = &ListNode{Data: node, Next: curNode.Next}
			lock.Unlock()
			return
		}
	}

	/* Case 4: Insertion at the Tail */
	curNode.Next = &ListNode{Data: node, Next: nil}
	lock.Unlock()
}

// process some of the operations that are queued up
func processQueue() {
	for {
		time.Sleep(2000 * time.Millisecond)
		lock.Lock()
		eliminateConcOps()
		processQueueHelper()
		lock.Unlock()
	}
}

// processQueueHelper does the actual processing of queue operations
func processQueueHelper() {
	updateCurTick()
	if queue != nil {
		eLog = eLog + "\n" + "BEFORE\n"
	}
	printQueue()

	for queue != nil {
		opNode := queue.Data

		/* Stop if any timestamp is exceeding the current safe tick */
		for i := 0; i < len(conns); i++ {
			if int(opNode.Timestamp["R"+strconv.Itoa(i+1)]) > curTick {
				if queue != nil {
					eLog = eLog + "\n" + "AFTER\n"
				}
				printQueue()
				return
			}
		}

		/* Run the associated op */
		if opNode.Type == IK {
			InsertKeyLocal(opNode.Key)
		} else if opNode.Type == IV {
			InsertValueLocal(opNode.Key, opNode.Value)
		}

		/* Remove Node */
		queue = queue.Next
	}
}

// updateCurTick updates the current "safe" tick
func updateCurTick() {
	noReplicas := len(conns)

	/* Gather all current timestamps */
	a := make([][]int, noReplicas)
	for n := queue; n != nil; n = n.Next {
		for i := 0; i < noReplicas; i++ {
			a[i] = append(a[i], int(n.Data.Timestamp["R"+strconv.Itoa(i+1)]))
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
	for i := 0; i < noReplicas; i++ {
		if len(a[i]) > 0 {
			eLog = eLog + fmt.Sprintln(a[i])
		}
	}
	eLog = eLog + ":" + fmt.Sprintln(curTick) + "\n"
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
