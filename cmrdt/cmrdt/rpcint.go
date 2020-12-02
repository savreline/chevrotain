package main

/* In this file
1. Definition and methods of RPC Int object
2. broadcastInsert method
*/

import (
	"fmt"
	"net/rpc"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"github.com/savreline/GoVector/govec/vclock"
)

// RPCIntArgs are the arguments all interval RPC Calls
type RPCIntArgs struct {
	Key, Value string
	Pid        string
	Clock      vclock.VClock
}

// InsertKeyRPC receives incoming insert key call
func (t *RPCInt) InsertKeyRPC(args *RPCIntArgs, reply *int) error {
	waitForTurn(args.Clock, args.Pid, args.Key, "", true)
	InsertKeyLocal(args.Key)
	return nil
}

// InsertValueRPC receives incoming insert value call
func (t *RPCInt) InsertValueRPC(args *RPCIntArgs, reply *int) error {
	waitForTurn(args.Clock, args.Pid, args.Key, args.Value, true)
	InsertValueLocal(args.Key, args.Value)
	return nil
}

// RemoveKeyRPC receives incoming insert key call
func (t *RPCInt) RemoveKeyRPC(args *RPCIntArgs, reply *int) error {
	waitForTurn(args.Clock, args.Pid, args.Key, "", false)
	RemoveKeyLocal(args.Key)
	return nil
}

// RemoveValueRPC receives incoming insert value call
func (t *RPCInt) RemoveValueRPC(args *RPCIntArgs, reply *int) error {
	waitForTurn(args.Clock, args.Pid, args.Key, args.Value, false)
	RemoveValueLocal(args.Key, args.Value)
	return nil
}

func broadcast(args *util.RPCExtArgs, insert bool) []*rpc.Call {
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
			if args.Value == "" && insert {
				fmt.Println("InsertKey RPC", no, "->", destNo)
				calls[i] = client.Go("RPCInt.InsertKeyRPC",
					RPCIntArgs{Key: args.Key, Pid: noStr, Clock: logger.GetCurrentVC()},
					&result, nil)
			} else if args.Value != "" && insert {
				fmt.Println("InsertValue RPC", no, "->", destNo)
				calls[i] = client.Go("RPCInt.InsertValueRPC",
					RPCIntArgs{Key: args.Key, Value: args.Value, Pid: noStr, Clock: logger.GetCurrentVC()},
					&result, nil)
			} else if args.Value == "" && !insert {
				fmt.Println("RemoveKey RPC", no, "->", destNo)
				calls[i] = client.Go("RPCInt.RemoveKeyRPC",
					RPCIntArgs{Key: args.Key, Pid: noStr, Clock: logger.GetCurrentVC()},
					&result, nil)
			} else {
				fmt.Println("RemoveValue RPC", no, "->", destNo)
				calls[i] = client.Go("RPCInt.RemoveValueRPC",
					RPCIntArgs{Key: args.Key, Value: args.Value, Pid: noStr, Clock: logger.GetCurrentVC()},
					&result, nil)
			}
			if err != nil {
				util.PrintErr(noStr, err)
			}
		}
	}

	return calls
}

// Ensure broadcast completes and (optionally) log error
func waitForCallsToComplete(calls []*rpc.Call) {
	for _, call := range calls {
		if call != nil {
			<-call.Done
		}
	}
}

// Wait for the correct turn for the incoming RPC call
func waitForTurn(incomingClock vclock.VClock, incomingPid string, key string, value string, insert bool) {
	/* Broadcast the incoming value to see if it helps any waiting call */
	broadcastClockValue(incomingClock)

	/* Check if this RPC call needs to wait */
	ready, cnt := logger.GetCurrentVCSafe().CompareBroadcastClock(incomingClock, incomingPid)
	eLog = eLog + fmt.Sprint("K: ", key) + fmt.Sprint(" count ", cnt) +
		fmt.Sprint(" Comparison: ", ready) + "\n"

	if ready == false {
		/* Make a channel to communicate on with this RPC call */
		channel := make(chan vclock.VClock, 10)

		/* Add the channel to the pool */
		lock.Lock()
		chans[channel] = channel
		lock.Unlock()

		/* Wait for the correct clock */
		i := 0
		for ok := true; ok != ready; i++ {
			recvClock := <-channel
			ready, cnt = logger.GetCurrentVCSafe().CompareBroadcastClock(recvClock, incomingPid)
			iLog = iLog + fmt.Sprint("K: ", key) + fmt.Sprint("count", cnt) +
				fmt.Sprint(" Comparison: ", ready) + fmt.Sprint(":", i) + "\n"
		}

		/* Remove channel from the pool */
		lock.Lock()
		delete(chans, channel)
		lock.Unlock()
	}

	/* Merge clock */
	var msg string
	if value == "" && insert {
		msg = "IN InsKey " + key + " from " + incomingPid
	} else if value != "" && insert {
		msg = "IN InsVal " + key + ":" + value + " from " + incomingPid
	} else if value == "" && !insert {
		msg = "IN RmvKey " + key + " from " + incomingPid
	} else {
		msg = "IN RmvVal " + key + ":" + value + " from " + incomingPid
	}
	logger.MergeIncomingClock(msg, incomingClock, govec.GetDefaultLogOptions().Priority)
}
