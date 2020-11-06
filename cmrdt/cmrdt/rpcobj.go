package cmrdt

/* In this file
Definition and methods of RPC Int object
*/

import (
	"fmt"

	"../../util"
)

// RPCInt is the RPC Object for internal replica-to-replica communication
type RPCInt int

// InsertKeyRPC receives incoming insert key call
func (t *RPCInt) InsertKeyRPC(args *KeyArgs, reply *int) error {
	util.PrintMsg(args.No, "Recv RPC Call to Insert "+args.Key)
	InsertKeyLocal(args.Key, args.No)
	return nil
}

// InsertValueRPC receives incoming insert value call
func (t *RPCInt) InsertValueRPC(args *ValueArgs, reply *int) error {
	fmt.Println("RPC Insert Value")
	return nil
}
