package main

/* In this file
1. Definition and methods of RPC Int object
2. broadcastInsert method
*/

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"time"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"github.com/savreline/GoVector/govec/vclock"
)

// KeyArgs are the arguments to the InsertKeyRPC call
type KeyArgs struct {
	Key   string
	Pid   string
	Clock vclock.VClock
}

// ValueArgs are the arguments to the InsertValueRPC call
type ValueArgs struct {
	Key, Value string
	Pid        string
	Clock      vclock.VClock
}

// InsertKeyRPC receives incoming insert key call
func (t *RPCInt) InsertKeyRPC(args *KeyArgs, reply *int) error {
	waitForTurn(args.Clock, args.Pid, args.Key, "")
	InsertKeyLocal(args.Key)
	return nil
}

// InsertValueRPC receives incoming insert value call
func (t *RPCInt) InsertValueRPC(args *ValueArgs, reply *int) error {
	waitForTurn(args.Clock, args.Pid, args.Key, args.Value)
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
					KeyArgs{Key: key, Pid: pid, Clock: logger.GetCurrentVC()},
					&result, nil)
			} else {
				fmt.Println("InsertValue RPC", no, "->", destNo)
				calls[i] = client.Go("RPCInt.InsertValueRPC",
					ValueArgs{Key: key, Value: value, Pid: pid, Clock: logger.GetCurrentVC()},
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
func waitForTurn(incomingClock vclock.VClock, incomingPid string, key string, value string) {
	/* Broadcast the incoming value to see if it helps any waiting call */
	broadcastClockValue(incomingClock)

	/* Check if this RPC call needs to wait */
	ready, str := logger.CompareBroadcastClock(incomingClock, incomingPid)
	eLog = eLog + fmt.Sprint("K: ", key) + str + fmt.Sprint(" Comparison: ", ready) + "\n"

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
			ready, str = logger.CompareBroadcastClock(recvClock, incomingPid)
			iLog = iLog + fmt.Sprint("K: ", key) + str + fmt.Sprint(" Comparison: ", ready) + fmt.Sprint(":", i) + "\n"
		}

		/* Remove channel from the pool */
		lock.Lock()
		delete(chans, channel)
		lock.Unlock()
	}

	/* Merge clock */
	var msg string
	if value == "" {
		msg = "IN InsKey " + key + " from " + incomingPid
	} else {
		msg = "IN InsVal " + key + ":" + value + " from " + incomingPid
	}
	logger.MergeIncomingClock(msg, incomingClock, govec.GetDefaultLogOptions().Priority)
}
