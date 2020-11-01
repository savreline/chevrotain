package cmrdt

import "fmt"

/**********************/
/*** 3: RPC METHODS ***/
/**********************/

// RPCObj is the RPC Object
type RPCObj int

// KeyArgs are the arguments to the InsertKeyRPC call
type KeyArgs struct {
	key string
}

// ValueArgs are the arguments to the InsertValueRPC call
type ValueArgs struct {
	key   string
	value string
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
