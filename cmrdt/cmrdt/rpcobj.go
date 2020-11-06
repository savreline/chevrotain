package main

/* In this file
Definition and methods of RPC Int object
*/

import (
	"fmt"

	"../../util"
)

// InsertKeyRPC receives incoming insert key call
func (t *RPCInt) InsertKeyRPC(args *util.KeyArgs, reply *int) error {
	InsertKeyLocal(args.Key)
	return nil
}

// InsertValueRPC receives incoming insert value call
func (t *RPCInt) InsertValueRPC(args *util.ValueArgs, reply *int) error {
	fmt.Println("RPC Insert Value")
	return nil
}
