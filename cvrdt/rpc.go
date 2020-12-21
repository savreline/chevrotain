package main

import (
	"fmt"
	"net/rpc"
	"time"

	"../util"
)

// StateArgs are the arguments to the MergeState RPC call
type StateArgs struct {
	PosState, NegState []util.DDoc
	SrcPid             string
	Timestamp          int
	Tick               int
}

// InsertKey inserts the given key into the positive collection
func (t *RPCExt) InsertKey(args *util.RPCExtArgs, reply *int) error {
	insertLocalRecord("Keys", args.Key, posCollection, nil)
	util.EmulateDelay(delay)
	return nil
}

// RemoveKey insert the given key into the negative collection
func (t *RPCExt) RemoveKey(args *util.RPCExtArgs, reply *int) error {
	insertLocalRecord("Keys", args.Key, negCollection, nil)
	util.EmulateDelay(delay)
	return nil
}

// InsertValue inserts value into the given key's record in the positive colleciton
func (t *RPCExt) InsertValue(args *util.RPCExtArgs, reply *int) error {
	insertLocalRecord(args.Key, args.Value, posCollection, nil)
	util.EmulateDelay(delay)
	return nil
}

// RemoveValue inserts value into the given key's record in the negative colleciton
func (t *RPCExt) RemoveValue(args *util.RPCExtArgs, reply *int) error {
	insertLocalRecord(args.Key, args.Value, negCollection, nil)
	util.EmulateDelay(delay)
	return nil
}

// MergeState merges incoming state with the current state at the replica
func (t *RPCInt) MergeState(args *StateArgs, reply *int) error {
	if verbose > 1 {
		util.PrintMsg(noStr, "Merging state from "+args.SrcPid)
	}

	/* Merge clock and state */
	lock.Lock()
	clock = util.Max(clock, args.Timestamp)
	mergeState(args.PosState, posCollection)
	mergeState(args.NegState, negCollection)

	/* With gc: If notified by the main replica about the current safe tick, accept it,
	merge collections, reply with own clock and broadcast own state. Without gc:
	merge collections at the end, when curSafeTick has been set to MAXTICK */
	if no != 1 && args.Tick != -1 {
		if gc {
			curSafeTick = args.Tick
			mergeCollections()
			*reply = mySafeTick
		} else if curSafeTick == MAXTICK {
			mergeCollections()
		}
		broadcast()
	}

	lock.Unlock()
	util.EmulateDelay(delay)
	return nil
}

// primary replica broadcasts state to all other replicas at intervals specified by timeInt
func runSU() {
	for {
		time.Sleep(time.Duration(timeInt) * time.Millisecond)
		if flagSU {
			lock.Lock()
			calls, results := broadcast()
			lock.Unlock()

			/* With gc: Merge collections and wait for the calls to complete
			to update the curSafeTick tick that way. Without gc: merge
			collections at the end, when curSafeTick has been set to MAXTICK */
			if gc {
				mergeCollections()
				curSafeTick = waitForBroadcastToFinish(calls, results)
			} else if curSafeTick == MAXTICK {
				mergeCollections()
			}
		}
	}
}

// broadcasts state to all other replicas
func broadcast() ([]*rpc.Call, []int) {
	var destNo int
	var flag = false
	var calls = make([]*rpc.Call, len(conns))
	var results = make([]int, len(conns))

	/* Tick the clock */
	clock++

	/* Fix own safe tick right prior to the broadcast,
	this way can only merge elements that been sent out */
	mySafeTick = clock

	/* If this is the main replica, broadcast the current safe tick */
	tick := -1
	if no == 1 {
		tick = curSafeTick
	}

	/* Download current state */
	posState := util.DownloadDState(db, "REPLICA "+noStr, posCollection, "0")
	negState := util.DownloadDState(db, "REPLICA "+noStr, negCollection, "0")
	state := StateArgs{
		PosState:  posState,
		NegState:  negState,
		SrcPid:    noStr,
		Timestamp: clock,
		Tick:      tick}

	if verbose > 0 {
		iLog = iLog + "\n\nSending State " + fmt.Sprint(mySafeTick) + ":" + fmt.Sprint(curSafeTick) + "\n"
		iLog = iLog + printDState(posState, "POSITIVE")
		iLog = iLog + printDState(negState, "NEGATIVE")
	}

	/* Broadcast */
	for i, client := range conns {
		if i == no {
			flag = true
		}
		if client != nil {
			if flag {
				destNo = i + 2
			} else {
				destNo = i + 1
			}
			if verbose > 1 {
				fmt.Println("RPC Merge State", no, "->", destNo)
			}
			calls[i] = client.Go("RPCInt.MergeState", state, &results[i], nil)
		}
	}

	return calls, results
}

// ensure broadcast completes
func waitForBroadcastToFinish(calls []*rpc.Call, results []int) int {
	result := mySafeTick
	for i, call := range calls {
		if call != nil {
			<-call.Done
			if results[i] < result {
				result = results[i]
			}
		}
	}
	return result
}
