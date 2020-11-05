package cvrdt

/* In this file
0. Definitions of StateArgs and RPCInt
1. MergeState RPC method
*/

// RPCInt is the RPC Object
type RPCInt int

// StateArgs are the arguments to the MergeState RPC
type StateArgs struct {
	No int
	// TODO
}

// MergeState is
func (t *RPCInt) MergeState(args *StateArgs, reply *int) error {
	return nil
}
