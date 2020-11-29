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
	NO
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
		return "Insert Key"
	} else if opCode == IV {
		return "Insert Value"
	} else if opCode == RK {
		return "Remove Key"
	} else if opCode == RV {
		return "Remove Value"
	} else if opCode == NO {
		return "No Op"
	} else {
		util.PrintErr(noStr, errors.New("lookupOpCode: unknown operation"))
		return ""
	}
}

// Print the Queue
func printQueue() {
	if queue != nil {
		eLog = eLog + "Queue\n"
	}
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
	if cmp == 2 {
		queue = &ListNode{Data: node, Next: queue}
		lock.Unlock()
		return
	}

	/* Case 3: Insertion Elsewhere: First Check that we can find a comparable node */
	curNode := queue
	for ; curNode.Next != nil; curNode = curNode.Next {
		cmpNext := node.Timestamp.CompareClocks(curNode.Next.Data.Timestamp)
		cmpPrev := node.Timestamp.CompareClocks(curNode.Data.Timestamp)

		if cmpNext == 2 {
			if cmpPrev == 1 {
				curNode.Data.ConcOp = true
			}
			curNode.Next = &ListNode{Data: node, Next: curNode.Next}
			lock.Unlock()
			return
		}
	}

	/* Case 4: Insertion at the Tail */
	cmp = node.Timestamp.CompareClocks(curNode.Data.Timestamp)
	if cmp == 3 {
		curNode.Next = &ListNode{Data: node, Next: curNode.Next}
		lock.Unlock()
		return
	}

	/* Case 5: Concurrent Insertion at the Head (gave up on finding a comparable node) */
	cmp = node.Timestamp.CompareClocks(queue.Data.Timestamp)
	if cmp == 1 {
		node.ConcOp = true
		if queue.Next != nil && node.Timestamp.CompareClocks(queue.Next.Data.Timestamp) == 1 {
			queue.Data.ConcOp = true
		}
		queue = &ListNode{Data: node, Next: queue}
		lock.Unlock()
		return
	}

	/* Case 6: Concurrent Insertion Elsewhere (gave up on finding a comparable node) */
	for curNode := queue; curNode != nil; curNode = curNode.Next {
		cmp := node.Timestamp.CompareClocks(curNode.Data.Timestamp)
		if cmp == 1 {
			curNode.Data.ConcOp = true
			if curNode.Next != nil && node.Timestamp.CompareClocks(curNode.Next.Data.Timestamp) == 1 {
				node.ConcOp = true
			}
			curNode.Next = &ListNode{Data: node, Next: curNode.Next}
			lock.Unlock()
			return
		}
	}
}

// process some of the operations that are queued up
func processQueue() {
	timeInt := <-channel
	close(channel)
	for {
		/* If nothing sent over last timeInt, send a no op */
		time.Sleep(time.Duration(timeInt) * time.Millisecond)
		if sent == false {
			processExtCall(util.RPCExtArgs{Key: "", Value: ""}, NO)
		}

		/* Process queue */
		lock.Lock()
		eLog = eLog + "\nAn Iteration\n"
		printQueue()
		processConcOps()
		printQueue()
		processQueueHelper()
		printQueue()
		sent = false
		lock.Unlock()
	}
}

// processQueueHelper does the actual processing of queue operations
func processQueueHelper() {
	updateCurTick()

	/* Reset ticks and prev pointer to be ready for further calls to processConcOps */
	ticks = make([][]int, noReplicas)
	prev = nil

	for queue != nil {
		opNode := queue.Data

		/* Stop if any timestamp is exceeding the current safe tick */
		for i := 0; i < noReplicas; i++ {
			if int(opNode.Timestamp["R"+strconv.Itoa(i+1)]) > curTick {
				return
			}
		}

		/* Run the associated op */
		if opNode.Type == IK {
			InsertKeyLocal(opNode.Key)
		} else if opNode.Type == IV {
			InsertValueLocal(opNode.Key, opNode.Value)
		} else if opNode.Type == RK {
			RemoveKeyLocal(opNode.Key)
		} else if opNode.Type == RV {
			RemoveValueLocal(opNode.Key, opNode.Value)
		}

		/* Remove Node */
		queue = queue.Next
	}
}

// updateCurTick updates the current "safe" tick
func updateCurTick() {
	if len(ticks[0]) == 0 {
		return
	}

	/* Determine the latest timestamp per replica */
	b := make([]int, noReplicas)
	for i := 0; i < noReplicas; i++ {
		b[i] = max(ticks, i)
	}

	/* Determine the earliest timestamp for all replicas */
	curTick = min(b)
	if curTick == 0 {
		curTick = 1
	}
	for i := 0; i < noReplicas; i++ {
		if len(ticks[i]) > 0 {
			eLog = eLog + fmt.Sprintln(ticks[i])
		}
	}
	if verbose == true {
		eLog = eLog + fmt.Sprintln(b) + "=======\n"
		eLog = eLog + ":" + fmt.Sprintln(curTick)
	}
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
			processBlock(first, n)

			/* Single-op block */
		} else if n.Data.ConcOp == false && inABlock == false {
			processSingleOpBlock(n)
		}
	}
}

// process a block of concurrent operations
func processBlock(first *ListNode, last *ListNode) {
	sameOp := checkIfSameOp(first)
	diffVals := checkIfDiffVals(first)
	if sameOp || diffVals {
		if prev == nil {
			queue = first
		} else {
			prev.Next = first
		}
		prev = last
		return
	}
	orderBlock(elimOps(first))
}

// process a single op block
func processSingleOpBlock(n *ListNode) {
	for i := 0; i < noReplicas; i++ {
		ticks[i] = append(ticks[i], int(n.Data.Timestamp["R"+strconv.Itoa(i+1)]))
	}
	if prev == nil {
		queue = n
	} else {
		prev.Next = n
	}
	prev = n
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
func checkIfDiffVals(first *ListNode) bool {
	setOfVals := make(map[string]bool)
	res := true

	for n := first; ; n = n.Next {

		/* Check for different values */
		val := n.Data.Key + ":" + n.Data.Value
		if setOfVals[val] {
			res = false
		}
		setOfVals[val] = true

		/* Add timestamps to ticks */
		for i := 0; i < noReplicas; i++ {
			ticks[i] = append(ticks[i], int(n.Data.Timestamp["R"+strconv.Itoa(i+1)]))
		}

		if n.Data.ConcOp == false {
			break
		}
	}

	return res
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
				switch settings[id-1] {
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

	return ops
}

// order the remaining operations in the block
func orderBlock(ops []map[string]int) {
	var flag = false
	var n *ListNode

	for _, id := range []int{1, 4, 2, 3} {
		for k := range ops[id] {
			args := strings.SplitN(k, ":", -1)
			opNode := OpNode{
				Type:      OpCode(id),
				Key:       args[0],
				Value:     args[1],
				Timestamp: make(map[string]uint64),
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
}
