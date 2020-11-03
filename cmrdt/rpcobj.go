package cmrdt

import "fmt"

/**********************/
/*** 3: RPC METHODS ***/
/**********************/

// RPCObj is the RPC Object
type RPCObj int

// ConnectArgs are the arguments to the ConnectReplica call (a dummy struct)
type ConnectArgs struct {
	Val string
}

// KeyArgs are the arguments to the InsertKeyRPC call
type KeyArgs struct {
	Key string
}

// ValueArgs are the arguments to the InsertValueRPC call
type ValueArgs struct {
	Key   string
	Value string
}

// InsertKeyRPC receives incoming insert key call
func (t *RPCObj) InsertKeyRPC(args *KeyArgs, reply *int) error {
	fmt.Println("RPC Insert Key")
	*reply = 200
	return nil
}

// InsertValueRPC receives incoming insert value call
func (t *RPCObj) InsertValueRPC(args *ValueArgs, reply *int) error {
	fmt.Println("RPC Insert Value")
	return nil
}
