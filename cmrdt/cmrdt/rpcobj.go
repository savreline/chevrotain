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

// KeyArgs are the arguments to the InsertKeyRPC call
type KeyArgs struct {
	Key   string
	Clock vclock.VClock
}

// ValueArgs are the arguments to the InsertValueRPC call
type ValueArgs struct {
	Key, Value string
	Clock      vclock.VClock
}

// InsertKeyRPC receives incoming insert key call
func (t *RPCInt) InsertKeyRPC(args *KeyArgs, reply *int) error {
	waitForTurn(args.Clock, args.Key, "")
	InsertKeyLocal(args.Key)
	return nil
}

// InsertValueRPC receives incoming insert value call
func (t *RPCInt) InsertValueRPC(args *ValueArgs, reply *int) error {
	waitForTurn(args.Clock, args.Key, args.Value)
	InsertValueLocal(args.Key, args.Value)
	return nil
}

func broadcastInsert(key string, value string) []*rpc.Call {
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
			if value == "" {
				fmt.Println("InsertKey RPC", no, "->", destNo)
				calls[i] = client.Go("RPCInt.InsertKeyRPC",
					KeyArgs{Key: key, Clock: logger.GetCurrentVC()},
					&result, nil)
			} else {
				fmt.Println("InsertValue RPC", no, "->", destNo)
				calls[i] = client.Go("RPCInt.InsertValueRPC",
					ValueArgs{Key: key, Value: value, Clock: logger.GetCurrentVC()},
					&result, nil)
			}
			if err != nil {
				util.PrintErr(err)
			}
		}
	}

	return calls
}

// Ensure broadcast completes and (optionally) log error
func waitForCallsToComplete(key string, value string, calls []*rpc.Call) {
	for _, call := range calls {
		if call != nil {
			<-call.Done
		}
	}
}

// Wait for the correct turn for the incoming RPC call
func waitForTurn(clock vclock.VClock, key string, value string) {
	/* Check if the RPC call needs to wait */
	wait := true // broadcastClockValue(logger.GetCurrentVC())

	if wait == true {
		/* Make a channel to communicate on with this RPC call */
		channel := make(chan vclock.VClock, 10)

		/* Add the channel to the pool */
		lock.Lock()
		chans[channel] = channel
		broadcastClockValue(clock) // to be moved
		lock.Unlock()

		/* Wait for the correct clock */
		<-channel

		/* Merge clock */
		var msg string
		if value == "" {
			msg = "IN InsKey " + key
		} else {
			msg = "IN InsVal " + key + ":" + value
		}
		logger.MergeIncomingClock(msg, clock, govec.GetDefaultLogOptions().Priority)

		/* Remove channel from the pool */
		lock.Lock()
		delete(chans, channel)
		lock.Unlock()
	}
}
