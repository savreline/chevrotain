package main

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"time"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"github.com/savreline/GoVector/govec/vclock"
)

// InsertKeyRPC receives incoming insert key call
func (t *RPCInt) InsertKeyRPC(args *KeyArgs, reply *int) error {
	queueCall(args.Key, "", args.Timestamp, args.Pid)
	InsertKeyLocal(args.Key)
	return nil
}

// InsertValueRPC receives incoming insert value call
func (t *RPCInt) InsertValueRPC(args *ValueArgs, reply *int) error {
	queueCall(args.Key, "", args.Timestamp, args.Pid)
	InsertValueLocal(args.Key, args.Value)
	return nil
}

// broadcastInsert
func broadcastInsert(key string, value string, timestamp vclock.VClock) []*rpc.Call {
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
			if value == "" {
				fmt.Println("InsertKey RPC", no, "->", destNo)
				calls[i] = client.Go("RPCInt.InsertKeyRPC",
					KeyArgs{Key: key, Timestamp: timestamp, Pid: noStr},
					&result, nil)
			} else {
				fmt.Println("InsertValue RPC", no, "->", destNo)
				calls[i] = client.Go("RPCInt.InsertValueRPC",
					ValueArgs{Key: key, Value: value, Timestamp: timestamp, Pid: noStr},
					&result, nil)
			}
			if err != nil {
				util.PrintErr(noStr, err)
			}
		}
	}
	return calls
}

// queueCall will place the call onto the queue
func queueCall(key string, value string, timestamp vclock.VClock, pid string) {
	/* Merge clock */
	var msg string
	if value == "" {
		msg = "IN InsKey " + key + " from " + pid
	} else {
		msg = "IN InsVal " + key + ":" + value + " from " + pid
	}
	logger.MergeIncomingClock(msg, timestamp, govec.GetDefaultLogOptions().Priority)
}

// Ensure broadcast completes
func waitForBroadcastToFinish(calls []*rpc.Call) {
	for _, call := range calls {
		if call != nil {
			<-call.Done
		}
	}
}
