package cmrdt

import (
	"fmt"

	"../../util"
)

/**********************/
/*** 3: RPC METHODS ***/
/**********************/

// RPCObj is the RPC Object
type RPCObj int

// ConnectArgs are the arguments to the ConnectReplica call (a dummy struct)
type ConnectArgs struct {
	No int
}

// KeyArgs are the arguments to the InsertKeyRPC call
type KeyArgs struct {
	No  int
	Key string
}

// ValueArgs are the arguments to the InsertValueRPC call
type ValueArgs struct {
	No         int
	Key, Value string
}

// InsertKeyRPC receives incoming insert key call
func (t *RPCObj) InsertKeyRPC(args *KeyArgs, reply *int) error {
	util.PrintMsg(args.No, "Recv RPC Call to Insert "+args.Key)
	InsertKeyLocal(args.Key, args.No)
	return nil
}

// InsertValueRPC receives incoming insert value call
func (t *RPCObj) InsertValueRPC(args *ValueArgs, reply *int) error {
	fmt.Println("RPC Insert Value")
	return nil
}
