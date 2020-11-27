package main

import (
	"fmt"
	"net/rpc"

	"../../util"
	"github.com/savreline/GoVector/govec"
)

// InsertKeyRPC receives incoming insert key call
func (t *RPCInt) InsertKeyRPC(args *OpNode, reply *int) error {
	processIntCall(*args)
	return nil
}

// InsertValueRPC receives incoming insert value call
func (t *RPCInt) InsertValueRPC(args *OpNode, reply *int) error {
	processIntCall(*args)
	return nil
}

// broadcast an operation
func broadcast(opNode OpNode) []*rpc.Call {
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
			if opNode.Value == "" {
				fmt.Println("InsertKey RPC", no, "->", destNo)
				calls[i] = client.Go("RPCInt.InsertKeyRPC", opNode, &result, nil)
			} else {
				fmt.Println("InsertValue RPC", no, "->", destNo)
				calls[i] = client.Go("RPCInt.InsertValueRPC", opNode, &result, nil)
			}
			if err != nil {
				util.PrintErr(noStr, err)
			}
		}
	}
	return calls
}

// process an external RPC call
func processExtCall(args util.RPCExtArgs, opCode OpCode) {
	logger.StartBroadcast("OUT "+noStr+": "+lookupOpCode(opCode)+" : "+args.Key+" : "+args.Value,
		govec.GetDefaultLogOptions())
	opNode := OpNode{
		Type:      opCode,
		Key:       args.Key,
		Value:     args.Value,
		Timestamp: logger.GetCurrentVC().Copy(),
		Pid:       noStr,
		ConcOp:    false}
	logger.StopBroadcast()
	addToQueue(opNode)
	calls := broadcast(opNode)
	waitForBroadcastToFinish(calls)
}

// process an internal RPC call
func processIntCall(opNode OpNode) {
	/* Add operation to queue */
	addToQueue(opNode)

	/* Merge clock */
	var msg string
	if opNode.Value == "" {
		msg = "IN InsKey " + opNode.Key + " from " + opNode.Pid
	} else {
		msg = "IN InsVal " + opNode.Key + ":" + opNode.Value + " from " + opNode.Pid
	}
	logger.MergeIncomingClock(msg, opNode.Timestamp, govec.GetDefaultLogOptions().Priority)
}

// Ensure broadcast completes
func waitForBroadcastToFinish(calls []*rpc.Call) {
	for _, call := range calls {
		if call != nil {
			<-call.Done
		}
	}
}
