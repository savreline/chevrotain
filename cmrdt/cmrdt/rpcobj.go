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
)

// InsertKeyRPC receives incoming insert key call
func (t *RPCInt) InsertKeyRPC(args *util.KeyArgs, reply *int) error {
	InsertKeyLocal(args.Key)
	return nil
}

// InsertValueRPC receives incoming insert value call
func (t *RPCInt) InsertValueRPC(args *util.ValueArgs, reply *int) error {
	InsertValueLocal(args.Key, args.Value)
	return nil
}

func broadcastInsert(key string, value string) {
	var result int
	var destNo int
	var err error
	var flag = false
	var calls = make([]*rpc.Call, len(conns))

	logger.StartBroadcast(govec.GetDefaultLogOptions())
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
				calls[i] = client.Go("RPCInt.InsertKeyRPC", util.KeyArgs{Key: key}, &result, nil)
			} else {
				fmt.Println("InsertValue RPC", no, "->", destNo)
				calls[i] = client.Go("RPCInt.InsertValueRPC", util.ValueArgs{Key: key, Value: value}, &result, nil)
			}
			if err != nil {
				util.PrintErr(err)
			}
		}
	}
	logger.StopBroadcast()

	// Ensure broadcast completes and (optionally) log error
	for i, call := range calls {
		if call != nil {
			replyCall := <-call.Done
			if verbose == true {
				eLog = eLog + fmt.Sprint("To Repl: ", i, " Key: ", key, " Val: ", value, " : ", replyCall.Error) + "\n"
			}
		}
	}
}
