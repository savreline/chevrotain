package cvrdt

/*In this file:
0. Definition of SnapShotArgs
1. GetCurrentSnapShot
*/

// SnapShotArgs are the arguments to the GetCurrentSnapShot RPC call
type SnapShotArgs struct {
	No int
}

// GetCurrentSnapShot merges the positive and negative datasets into one collection
func (t *RPCExt) GetCurrentSnapShot(args *SnapShotArgs, reply *int) error {
	// TODO
	return nil
}
