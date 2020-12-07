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
	if verbose {
		util.PrintMsg(noStr, "Merging state from "+args.SrcPid)
	}

	/* Merge clock and state */
	clock = util.Max(clock, args.Timestamp)
	mergeState(args.PosState, posCollection)
	mergeState(args.NegState, negCollection)

	/* If notified by the main replica about the current safe tick, accept it
	and reply with own clock (this will never run on the main replicas as it
	cannot notify itself) */
	if args.Tick != -1 {
		curSafeTick = args.Tick
		mergeCollections()
		*reply = clock
	}
	util.EmulateDelay(delay)
	return nil
}

// broadcasts state to all other replicas at intervals specified by timeInt
func runSU() {
	<-chanSU
	for {
		time.Sleep(time.Duration(timeInt) * time.Millisecond)
		calls, results := broadcast()

		/* If this is the main replica, wait for the calls to complete
		and update the curSafeTick tick that way */
		if no == 1 {
			curSafeTick = waitForBroadcastToFinish(calls, results)
			mergeCollections()
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

	/* If the main replica, broadcast the current safe tick */
	tick := -1
	if no == 1 {
		tick = curSafeTick
	}

	/* Download current state */
	posState := util.DownloadDState(db.Collection(posCollection), "REPLICA "+noStr, "0")
	negState := util.DownloadDState(db.Collection(negCollection), "REPLICA "+noStr, "0")
	state := StateArgs{
		PosState:  posState,
		NegState:  negState,
		SrcPid:    noStr,
		Timestamp: clock,
		Tick:      tick}

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
			if verbose {
				fmt.Println("RPC Merge State", no, "->", destNo)
			}
			calls[i] = client.Go("RPCInt.MergeState", state, &results[i], nil)
		}
	}
	return calls, results
}

// ensure broadcast completes
func waitForBroadcastToFinish(calls []*rpc.Call, results []int) int {
	result := -1
	for i, call := range calls {
		if call != nil {
			<-call.Done
			if results[i] > result {
				result = results[i]
			}
		}
	}
	return result
}
