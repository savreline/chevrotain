package main

import (
	"fmt"
	"net/rpc"
	"time"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"github.com/savreline/GoVector/govec/vclock"
)

// StateArgs are the arguments to the MergeState RPC
type StateArgs struct {
	PosState, NegState []util.CvRecord
	No                 string
	Timestamp          vclock.VClock
}

// MergeState merges incoming state with the current state at the replica
func (t *RPCInt) MergeState(args *StateArgs, reply *int) error {
	msg := "Recv State from " + args.No
	logger.MergeIncomingClock(msg, args.Timestamp, govec.GetDefaultLogOptions().Priority)
	mergeCollection(args.PosState, posCollection)
	mergeCollection(args.PosState, negCollection)
	return nil
}

// GetCurrentSnapShot merges the positive and negative datasets into one collection
func (t *RPCExt) GetCurrentSnapShot(args *util.RPCExtArgs, reply *int) error {
	// TODO
	return nil
}

// broadcasts state to all other replicas
func sendState() {
	for {
		time.Sleep(time.Duration(timeInt) * time.Millisecond)
		posState := util.DownloadCvState(db.Collection(posCollection), "0")
		negState := util.DownloadCvState(db.Collection(negCollection), "0")
		logger.StartBroadcast("OUT "+noStr, govec.GetDefaultLogOptions())
		broadcast(StateArgs{PosState: posState,
			NegState:  negState,
			No:        noStr,
			Timestamp: logger.GetCurrentVC().Copy()})
		logger.StopBroadcast()
	}
}

// broadcasts state to all other replicas
func broadcast(state StateArgs) []*rpc.Call {
	var result int
	var destNo int
	var err error
	var flag = false
	var calls = make([]*rpc.Call, len(conns))

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
			fmt.Println("RPC Sending State", no, "->", destNo)
			calls[i] = client.Go("RPCInt.MergeState", state, &result, nil)
			if err != nil {
				util.PrintErr(noStr, err)
			}
		}
	}
	return calls
}
