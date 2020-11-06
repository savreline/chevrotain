package cmrdt

/* In this file
0. Definitions of ValueArgs
1. InsertValue Ext RPC method
2. InsertValueLocal (that works with the local db), InsertValueGlobal (that broadcats the event) methods
*/

import (
	"context"
	"fmt"

	"../../util"
	"github.com/DistributedClocks/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

// ValueArgs are the arguments to the InsertValueRPC call
type ValueArgs struct {
	No         int
	Key, Value string
}

// InsertValue inserts value into the given key
func (t *RPCExt) InsertValue(args *ValueArgs, reply *int) error {
	InsertValueLocal(args.Key, args.Value, args.No)
	InsertValueGlobal(args.Key, args.Value, args.No)
	return nil
}

// InsertValueLocal inserts the value into the local db
func InsertValueLocal(key string, value string, no int) {
	replicas[no].logger.LogLocalEvent("Inserting value"+value, govec.GetDefaultLogOptions())
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: value}}}}

	updateResult, err := replicas[no].db.Collection("kvs").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)
}

// InsertValueGlobal broadcasts the insertValue operation to other replicas
func InsertValueGlobal(key string, value string, no int) {
	var result int
	err := replicas[no].clients[0].Call("InsertValueRPC", ValueArgs{no, key, value}, &result)
	if err != nil {
		util.PrintErr(err)
	}
}
