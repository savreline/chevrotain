package cvrdt

/* In this file
0. Definitions of KeyArgs
1. InsertKey, RemoveKey RPC methods
2. InsertKeyHelper local method that works with db
*/

import (
	"context"
	"fmt"
	"strconv"

	"../../util"
	"github.com/DistributedClocks/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

// KeyArgs are the arguments to the InsertKeyRPC call
type KeyArgs struct {
	No  int
	Key string
}

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *KeyArgs, reply *int) error {
	InsertKeyHelper(args.Key, args.No, posCollection)
	return nil
}

// RemoveKey removes the given key
func (t *RPCExt) RemoveKey(args *KeyArgs, reply *int) error {
	InsertKeyHelper(args.Key, args.No, negCollection)
	return nil
}

// InsertKeyHelper inserts the key in either positive collection (add) or negative collection (remove)
func InsertKeyHelper(key string, no int, collection string) {
	logger := replicas[no].logger
	db := replicas[no].db
	logger.LogLocalEvent("Inserting Key "+key, govec.GetDefaultLogOptions())

	/* Update global keys entry */
	valueEntry := ValueEntry{key, replicas[no].logger.GetCurrentVC()}
	filter := bson.D{{Key: "name", Value: "Keys"}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: valueEntry}}}}
	updateResult, err := db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Printf("REPLICA "+strconv.Itoa(no+1)+": Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)

	/* Insert the Key */
	newRecord := Record{key, logger.GetCurrentVC(), []ValueEntry{}}
	_, err = db.Collection(collection).InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}
	util.PrintMsg(no, "Inserted Key "+key)
}
