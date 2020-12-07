package main

import (
	"fmt"
	"net/rpc"

	"../util"
)

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *util.RPCExtArgs, reply *int) error {
	calls := broadcast(args, false)
	waitForBroadcastToFinish(calls)
	util.EmulateDelay(delay)
	return nil
}

// InsertValue inserts value into the given key
func (t *RPCExt) InsertValue(args *util.RPCExtArgs, reply *int) error {
	calls := broadcast(args, false)
	waitForBroadcastToFinish(calls)
	util.EmulateDelay(delay)
	return nil
}

// RemoveKey removes the given key
func (t *RPCExt) RemoveKey(args *util.RPCExtArgs, reply *int) error {
	calls := broadcast(args, true)
	waitForBroadcastToFinish(calls)
	util.EmulateDelay(delay)
	return nil
}

// RemoveValue removes the given value from the given key
func (t *RPCExt) RemoveValue(args *util.RPCExtArgs, reply *int) error {
	calls := broadcast(args, true)
	waitForBroadcastToFinish(calls)
	util.EmulateDelay(delay)
	return nil
}

// broadcast an operation
func broadcast(args *util.RPCExtArgs, remove bool) []*rpc.Call {
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
			if args.Value == "" && !remove {
				calls[i] = client.Go("RPCInt.InsertKey", args, &result, nil)
			} else if args.Value != "" && !remove {
				calls[i] = client.Go("RPCInt.RemoveKey", args, &result, nil)
			} else if args.Value == "" && remove {
				calls[i] = client.Go("RPCInt.InsertValue", args, &result, nil)
			} else {
				calls[i] = client.Go("RPCInt.RemoveValue", args, &result, nil)
			}
		}
	}
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
