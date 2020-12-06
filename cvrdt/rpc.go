package main

import (
	"fmt"
	"time"

	"../util"
)

// StateArgs are the arguments to the MergeState RPC call
type StateArgs struct {
	PosState, NegState []util.DDoc
	SrcPid             string
	Timestamp          int
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
	clock = util.Max(clock, args.Timestamp)
	mergeState(args.PosState, posCollection)
	mergeState(args.NegState, negCollection)
	util.EmulateDelay(delay)
	return nil
}

// broadcasts state to all other replicas at intervals specified by timeInt
func runSU() {
	<-chanSU
	for {
		time.Sleep(time.Duration(timeInt) * time.Millisecond)
		broadcast()
	}
}

// runs garbage collection at intervals specified by timeInt
func runGC() {
	<-chanGC
	for {
		time.Sleep(time.Duration(timeInt) * time.Millisecond)
		mergeCollections()
	}
}

// broadcasts state to all other replicas
func broadcast() {
	var result int
	var destNo int
	var flag = false

	/* Tick the clock */
	clock++

	/* Download current state */
	posState := util.DownloadDState(db.Collection(posCollection), "REPLICA "+noStr, "0")
	negState := util.DownloadDState(db.Collection(negCollection), "REPLICA "+noStr, "0")
	state := StateArgs{
		PosState:  posState,
		NegState:  negState,
		SrcPid:    noStr,
		Timestamp: clock}

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
			client.Go("RPCInt.MergeState", state, &result, nil)
		}
	}
}
