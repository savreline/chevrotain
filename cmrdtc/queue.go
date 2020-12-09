package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"../util"
)

// Global variables: the prev pointer refers to the immediately preceeding block
// of the queue when re-linking the queue blocks while processing the queue
var prev *ListNode

// ListNode is a node in the linked list queue with data and a next pointer
type ListNode struct {
	Data OpNode
	Next *ListNode
}

// print the queue to the log
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
	var curNode *ListNode
	var cmp int

	/* Case 1: Empty Queue */
	if queue == nil {
		queue = &ListNode{Data: node, Next: nil}
		goto finish
	}

	/* Case 2: Insertion at the Head */
	cmp = node.Timestamp.CompareClocks(queue.Data.Timestamp)
	if cmp == 2 {
		queue = &ListNode{Data: node, Next: queue}
		goto finish
	}

	/* Case 3: Insertion Elsewhere: First Check that we can find a comparable node */
	curNode = queue
	for ; curNode.Next != nil; curNode = curNode.Next {
		cmpNext := node.Timestamp.CompareClocks(curNode.Next.Data.Timestamp)
		cmpPrev := node.Timestamp.CompareClocks(curNode.Data.Timestamp)

		if cmpNext == 2 {
			if cmpPrev == 1 {
				curNode.Data.ConcOp = true
			}
			curNode.Next = &ListNode{Data: node, Next: curNode.Next}
			goto finish
		}
	}

	/* Case 4: Insertion at the Tail */
	cmp = node.Timestamp.CompareClocks(curNode.Data.Timestamp)
	if cmp == 3 {
		curNode.Next = &ListNode{Data: node, Next: curNode.Next}
		goto finish
	}

	/* Case 5: Concurrent Insertion at the Head (gave up on finding a comparable node) */
	cmp = node.Timestamp.CompareClocks(queue.Data.Timestamp)
	if cmp == 1 {
		node.ConcOp = true
		if queue.Next != nil && node.Timestamp.CompareClocks(queue.Next.Data.Timestamp) == 1 {
			queue.Data.ConcOp = true
		}
		queue = &ListNode{Data: node, Next: queue}
		goto finish
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
			goto finish
		}
	}

finish:
	queueLen++
	lock.Unlock()
	if queueLen > MAXQUEUELEN {
		processQueue()
	}
}

// process some of the operations that are queued up
func processQueue() {
	eLog = eLog + "\nAn Iteration\n"
	lock.Lock()
	printQueue()
	processConcOps()
	printQueue()
	processQueueHelper()
	printQueue()
	lock.Unlock()
}

// this method does the actual processing of queue operations
func processQueueHelper() {
	updateCurTick()

	/* Reset ticks and prev pointer to be ready for further calls to processConcOps */
	ticks = make([][]int, noReplicas)
	prev = nil

	for queue != nil {
		opNode := queue.Data

		/* Stop if any timestamp is exceeding the current safe tick */
		for i := 0; i < noReplicas; i++ {
			if int(opNode.Timestamp["R"+strconv.Itoa(i+1)]) > curSafeTick {
				return
			}
		}

		/* Run the associated op */
		if opNode.OpType == util.IK {
			insertKey(opNode.Key)
		} else if opNode.OpType == util.IV {
			insertValue(opNode.Key, opNode.Value)
		} else if opNode.OpType == util.RK {
			removeKey(opNode.Key)
		} else if opNode.OpType == util.RV {
			removeValue(opNode.Key, opNode.Value)
		}

		/* Remove Node */
		queue = queue.Next
		queueLen--
	}
}

// updates the current "safe" tick
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
	curSafeTick = min(b)
	if curSafeTick == 0 {
		curSafeTick = 1
	}

	/* Add information about min picks to the log */
	for i := 0; i < noReplicas; i++ {
		if verbose && len(ticks[i]) > 0 {
			eLog = eLog + fmt.Sprintln(ticks[i])
		}
	}
	if verbose {
		eLog = eLog + fmt.Sprintln(b) + "=======\n"
		eLog = eLog + ":" + fmt.Sprintln(curSafeTick)
	}
}

// determines the latest timestamp per replica
func max(a [][]int, i int) int {
	sort.Ints(a[i])
	for j := 0; j < len(a[i])-1; j++ {
		if a[i][j+1] > curSafeTick+1 && a[i][j+1] > a[i][j]+1 {
			return a[i][j]
		}
	}
	return a[i][len(a[i])-1]
}

// determines the earliest timestamp for all replicas
func min(b []int) int {
	res := b[0]
	for j := 1; j < len(b); j++ {
		if b[j] < res {
			res = b[j]
		}
	}
	return res
}

// process concurrent operations from the queue using the predefined preference
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
	val := first.Data.OpType
	for n := first; ; n = n.Next {
		if n.Data.OpType != val {
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
	/* Init maps: maps an operation type to the number of times this operation
	has been encountered in the current block */
	ops := make([]map[string]int, 10)
	for i := 1; i < 10; i++ {
		ops[i] = make(map[string]int)
	}

	/* Build up counts */
	for n := first; ; n = n.Next {
		val := n.Data.Key + ":" + n.Data.Value
		ops[n.Data.OpType][val]++
		if n.Data.ConcOp == false {
			break
		}
	}

	/* Eliminate as per bias: other loop loops around key (=1) and val (=2) ops.
	The inner loop loops aroud ops with a specific key and value pair */
	for _, id := range []int{1, 2} {
		for k := range ops[id] {
			if ops[id+2][k] > 0 { // if exist removes for this operation
				switch bias[id-1] {
				case true:
					ops[id][k] = 0   // adds
					ops[id+2][k] = 1 // removes
				case false:
					ops[id][k] = 1   // adds
					ops[id+2][k] = 0 // removes
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
				OpType:    util.OpCode(id),
				Key:       args[0],
				Value:     args[1],
				Timestamp: make(map[string]uint64),
				SrcPid:    "0",
				ConcOp:    false}

			/* If this is the first op in a block, must link previous block */
			if !flag {
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
