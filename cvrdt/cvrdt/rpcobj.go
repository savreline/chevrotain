package cvrdt

// RPCObj is the RPC Object
type RPCObj int

// StateArgs are the arguments to the MergeState RPC
type StateArgs struct {
	No int
	// TODO
}

// MergeState is
func (t *RPCObj) MergeState(args *StateArgs, reply *int) error {
	return nil
}
