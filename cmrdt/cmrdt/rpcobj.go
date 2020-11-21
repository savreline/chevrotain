package main

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"time"

	"../../util"
	"github.com/savreline/GoVector/govec"
)

// InsertKeyRPC receives incoming insert key call
func (t *RPCInt) InsertKeyRPC(args *OpNode, reply *int) error {
	queueCall(*args)
	return nil
}

// InsertValueRPC receives incoming insert value call
func (t *RPCInt) InsertValueRPC(args *OpNode, reply *int) error {
	queueCall(*args)
	return nil
}

// broadcastInsert
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
		if delay > 0 {
			time.Sleep(time.Duration(rand.Intn(delay)) * time.Millisecond)
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

// queueCall will place the call onto the queue
func queueCall(opNode OpNode) {
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

// Make a copy of the current clock
func copyCurrentClock() map[string]uint64 {
	timestamp := make(map[string]uint64, len(logger.GetCurrentVC()))
	for k, v := range logger.GetCurrentVC() {
		timestamp[k] = v
	}
	return timestamp
}
