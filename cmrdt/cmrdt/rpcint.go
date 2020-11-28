package main

import (
	"fmt"
	"net/rpc"

	"../../util"
	"github.com/savreline/GoVector/govec"
)

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
			if verbose == true {
				fmt.Println(lookupOpCode(opNode.Type)+" RPC", no, "->", destNo)
			}
			calls[i] = client.Go("RPCInt.ProcessIntCall", opNode, &result, nil)
			if err != nil {
				util.PrintErr(noStr, err)
			}
		}
	}
	sent = true
	return calls
}

// process an external RPC call:
//	1. tick the clock
//	2. package the operation into an opNode struct
//	3. add to local queue
//	4. call the broadcast method
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

// ProcessIntCall processes an internal RPC call:
//	1. add to local queue
//	2. merge the incoming clock
func (t *RPCInt) ProcessIntCall(args *OpNode, reply *int) error {
	addToQueue(*args)

	var msg string
	if args.Value == "" {
		msg = "IN InsKey " + args.Key + " from " + args.Pid
	} else {
		msg = "IN InsVal " + args.Key + ":" + args.Value + " from " + args.Pid
	}
	logger.MergeIncomingClock(msg, args.Timestamp, govec.GetDefaultLogOptions().Priority)
	return nil
}

// Ensure broadcast completes
func waitForBroadcastToFinish(calls []*rpc.Call) {
	for _, call := range calls {
		if call != nil {
			<-call.Done
		}
	}
}
