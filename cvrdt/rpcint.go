package main

import (
	"fmt"
	"time"

	"../util"
)

// StateArgs are the arguments to the MergeState RPC
type StateArgs struct {
	PosState, NegState []util.CvDoc
	SrcPid             string
	Timestamp          int
}

// MergeState merges incoming state with the current state at the replica
func (t *RPCInt) MergeState(args *StateArgs, reply *int) error {
	if verbose {
		util.PrintMsg(noStr, "Merging state from "+args.SrcPid)
	}
	clock = util.Max(clock, args.Timestamp)
	mergeState(args.PosState, posCollection)
	mergeState(args.NegState, negCollection)
	return nil
}

// broadcasts state to all other replicas
func runSE() {
	<-chanSE
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
		gc()
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
	posState := util.DownloadCvState(db.Collection(posCollection), "REPLICA "+noStr, "0")
	negState := util.DownloadCvState(db.Collection(negCollection), "REPLICA "+noStr, "0")
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

func gc() {
}
