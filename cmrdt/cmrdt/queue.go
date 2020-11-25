package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"../../util"
	"github.com/savreline/GoVector/govec/vclock"
)

// Global variables
var prev *ListNode

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
	eLog = eLog + "Queue" + "\n"
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
		processConcOps()
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

// process concurrent operations from the queue using predefined preference
func processConcOps() {
	var inABlock = false
	var first *ListNode
	for n := queue; n != nil; n = n.Next {

		/* Start of a block */
		if n.Data.ConcOp == true && inABlock == false {
			inABlock = true
			first = n

			/* End of a block */
		} else if n.Data.ConcOp == false && inABlock == true {
			inABlock = false
			processBlock(first)

			/* Single-op block */
		} else if n.Data.ConcOp == false && inABlock == false {
			prev.Next = n
			prev = n
		}
	}
}

// process a block of concurrent operations
func processBlock(first *ListNode) {
	sameOp := checkIfSameOp(first)
	diffVals, maxClock := checkIfDiffVals(first)
	if sameOp || diffVals {
		// return
	}
	ops := elimOps(first)
	orderBlock(maxClock, ops)
}

// shortcut no. 1: all operations are the same
func checkIfSameOp(first *ListNode) bool {
	val := first.Data.Type
	for n := first; ; n = n.Next {
		if n.Data.Type != val {
			return false
		}
		if n.Data.ConcOp == false {
			break
		}
	}
	return true
}

// shortcut no. 2: all values are different
func checkIfDiffVals(first *ListNode) (bool, vclock.VClock) {
	noReplicas := len(conns)
	a := make([][]int, noReplicas)
	setOfVals := make(map[string]bool)
	res := true

	/* Check for different values and add timestamps to a */
	for n := first; ; n = n.Next {
		val := n.Data.Key + ":" + n.Data.Value
		if setOfVals[val] {
			res = false
		}
		setOfVals[val] = true
		for i := 0; i < noReplicas; i++ {
			a[i] = append(a[i], int(n.Data.Timestamp["R"+strconv.Itoa(i+1)]))
		}
		if n.Data.ConcOp == false {
			break
		}
	}

	/* Determine b */
	b := make(map[string]uint64, noReplicas)
	for i := 0; i < noReplicas; i++ {
		b["R"+strconv.Itoa(i+1)] = uint64(max(a, i))
	}

	return res, b
}

// count number of operations of each type and eliminate accordinly
func elimOps(first *ListNode) []map[string]int {
	/* Init maps */
	ops := make([]map[string]int, 5)
	for i := 1; i < 5; i++ {
		ops[i] = make(map[string]int)
	}

	/* Build up counts */
	for n := first; ; n = n.Next {
		val := n.Data.Key + ":" + n.Data.Value
		ops[n.Data.Type][val]++
		if n.Data.ConcOp == false {
			break
		}
	}

	/* Eliminate as per spec */
	for _, id := range []int{1, 2} {
		for k := range ops[id] {
			if ops[id+2][k] > 0 {
				switch flag[id-1] {
				case -1:
					ops[id][k] = 0
					ops[id+2][k] = 1
				case 0:
					ops[id][k] = 0
					ops[id+2][k] = 0
				case 1:
					ops[id][k] = 1
					ops[id+2][k] = 0
				}
			}
		}
	}

	/* Print for Testing */
	// for i := 1; i < 5; i++ {
	// 	fmt.Println(i)
	// 	for k, v := range ops[i] {
	// 		fmt.Println(k, ":", v)
	// 	}
	// 	fmt.Println()
	// }
	return ops
}

// order the remaining operations in the block
func orderBlock(timestamp vclock.VClock, ops []map[string]int) {
	var flag = false
	var n *ListNode

	for _, id := range []int{1, 4, 2, 3} {
		for k := range ops[id] {
			args := strings.SplitN(k, ":", -1)
			opNode := OpNode{
				Type:      OpCode(id),
				Key:       args[0],
				Value:     args[1],
				Timestamp: timestamp,
				Pid:       "0",
				ConcOp:    false}

			/* If this is the first op in a block, must link previous block */
			if flag == false {
				flag = true

				/* No prev, link the queue pointer */
				if prev == nil {
					queue = &ListNode{Data: opNode, Next: nil}
					n = queue

					/* Else, link prev */
				} else {
					prev.Next = &ListNode{Data: opNode, Next: nil}
					n = prev.Next

				}

				continue
			}

			/* Linking with the block */
			n.Next = &ListNode{Data: opNode, Next: nil}
			n = n.Next
		}
	}

	/* Set prev to point to the last node */
	prev = n
	printQueue()
}
