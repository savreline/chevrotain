package main

import (
	"fmt"
	"net/rpc"
	"time"

	"../util"
	"github.com/savreline/GoVector/govec"
	"github.com/savreline/GoVector/govec/vclock"
)

// OpNode contains all information about a particular operations
// and represents the operation in the queue
type OpNode struct {
	Type       util.OpCode
	Key, Value string
	Timestamp  vclock.VClock
	Pid        string
	ConcOp     bool
}

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *util.RPCExtArgs, reply *int) error {
	processExtCall(*args, util.IK)
	return nil
}

// InsertValue inserts value into the given key
func (t *RPCExt) InsertValue(args *util.RPCExtArgs, reply *int) error {
	processExtCall(*args, util.IV)
	return nil
}

// RemoveKey removes the given key
func (t *RPCExt) RemoveKey(args *util.RPCExtArgs, reply *int) error {
	processExtCall(*args, util.RK)
	return nil
}

// RemoveValue removes the given value from the given key
func (t *RPCExt) RemoveValue(args *util.RPCExtArgs, reply *int) error {
	processExtCall(*args, util.RV)
	return nil
}

// process an external RPC call
func processExtCall(args util.RPCExtArgs, opCode util.OpCode) {
	/* Tick the clock */
	logger.StartBroadcast("OUT "+noStr+": "+util.LookupOpCode(opCode, noStr)+" : "+args.Key+" : "+args.Value,
		govec.GetDefaultLogOptions())

	/* Package the operation into an OpNode struct */
	opNode := OpNode{
		Type:      opCode,
		Key:       args.Key,
		Value:     args.Value,
		Timestamp: logger.GetCurrentVC().Copy(),
		Pid:       noStr,
		ConcOp:    false}
	logger.StopBroadcast() // timestamp saved, ok to release the clock lock

	/* Add the operation to the local queue */
	addToQueue(opNode)

	/* Do the broadcast */
	calls := broadcast(opNode)
	waitForBroadcastToFinish(calls)
	util.EmulateDelay(delay)
}

// periodically broadcast no ops to keep clocks reasonably close
func runNoOps() {
	time.Sleep(time.Duration(timeInt) * time.Millisecond)
	if !sent {
		processExtCall(util.RPCExtArgs{Key: "", Value: ""}, util.NO)
	}
	sent = false // start a new time interval
}

// broadcast an operation
func broadcast(opNode OpNode) []*rpc.Call {
	var result int
	var destNo int
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
			if verbose {
				fmt.Println("RPC Int", no, "->", destNo)
			}
			calls[i] = client.Go("RPCInt.ProcessIntCall", opNode, &result, nil)
		}
	}
	sent = true
	return calls
}

// ensure broadcast completes
func waitForBroadcastToFinish(calls []*rpc.Call) {
	for _, call := range calls {
		if call != nil {
			<-call.Done
		}
	}
}
