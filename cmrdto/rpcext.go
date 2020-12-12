package main

import (
	"fmt"
	"net/rpc"

	"../util"
	"github.com/savreline/GoVector/govec"
	"github.com/savreline/GoVector/govec/vclock"
)

// BroadcastArgs are the arguments to all internal RPC calls
type BroadcastArgs struct {
	OpType     util.OpCode
	Key, Value string
	SrcPid     string
	Clock      vclock.VClock
	Ids        []int
}

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *util.RPCExtArgs, reply *int) error {
	processExtCall(args.Key, args.Value, util.IK)
	return nil
}

// InsertValue inserts value into the given key
func (t *RPCExt) InsertValue(args *util.RPCExtArgs, reply *int) error {
	processExtCall(args.Key, args.Value, util.IV)
	return nil
}

// RemoveKey removes the given key
func (t *RPCExt) RemoveKey(args *util.RPCExtArgs, reply *int) error {
	processExtCall(args.Key, args.Value, util.RK)
	return nil
}

// RemoveValue removes value from the given key
func (t *RPCExt) RemoveValue(args *util.RPCExtArgs, reply *int) error {
	processExtCall(args.Key, args.Value, util.RV)
	return nil
}

func processExtCall(key string, value string, opType util.OpCode) {
	/* Get and current clock and then tick the clock (and acquire the lock in doing so) */
	msg := util.LookupOpCode(opType, noStr) + " on R" + noStr + " K:" + key + " V:" + value
	clock := logger.StartBroadcast(msg, govec.GetDefaultLogOptions())

	/* Prepare-updates */
	var ids []int
	if opType == util.RK {
		ids = prepareUpdates("Keys", key, opType)
	} else {
		ids = prepareUpdates(key, value, opType)
	}

	/* Broadcast */
	calls := broadcast(BroadcastArgs{
		OpType: opType,
		Key:    key,
		Value:  value,
		SrcPid: noStr,
		Ids:    ids,
		Clock:  clock,
	})

	/* Effect-updates */
	effectUpdates(key, value, ids, opType)

	/* Release the lock as effect updates are now complete */
	logger.StopBroadcast()
	waitForBroadcastToFinish(calls)
	util.EmulateDelay(delay)
}

// broadcasts state to all other replicas
func broadcast(args BroadcastArgs) []*rpc.Call {
	var result int
	var destNo int
	var flag = false
	var calls = make([]*rpc.Call, len(conns))

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
				fmt.Println("RPC Int", no, "->", destNo)
			}
			calls[i] = client.Go("RPCInt.ProcessIntCall", args, &result, nil)
		}
	}
	return calls
}

// prepapre-updates
func prepareUpdates(key string, value string, opType util.OpCode) []int {
	if opType == util.IK || opType == util.IV {
		id++
		return []int{id}
	} else if opType == util.RK {
		return computeRemovalSet("Keys", value)
	}
	return computeRemovalSet(key, value)
}

// effect-updates
func effectUpdates(key string, value string, ids []int, opType util.OpCode) {
	if opType == util.IK {
		insert("Keys", key, ids[0])
	} else if opType == util.IV {
		insert(key, value, ids[0])
	} else if opType == util.RK {
		remove("Keys", key, ids)
	} else {
		remove(key, value, ids)
	}
}

// ensure broadcast completes
func waitForBroadcastToFinish(calls []*rpc.Call) {
	for _, call := range calls {
		if call != nil {
			<-call.Done
		}
	}
}
