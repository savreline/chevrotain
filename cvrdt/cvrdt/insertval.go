package cvrdt

/* In this file
0. Definitions of ValueArgs
1. InsertValue, RemoveValue RPC methods
2. InsertValueHelper local method that works with db
*/

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

// ValueArgs are the arguments to the InsertValueRPC call
type ValueArgs struct {
	No         int
	Key, Value string
}

// InsertValue inserts value into the given key
func (t *RPCExt) InsertValue(args *ValueArgs, reply *int) error {
	InsertValueHelper(args.Key, args.Value, args.No, posCollection)
	return nil
}

// RemoveValue removes value from the given key
func (t *RPCExt) RemoveValue(args *ValueArgs, reply *int) error {
	InsertValueHelper(args.Key, args.Value, args.No, negCollection)
	return nil
}

// InsertValueHelper inserts the value in either positive collection (add) or negative collection (remove)
func InsertValueHelper(key string, value string, no int, collection string) {
	logger := replicas[no].logger
	db := replicas[no].db
	logger.LogLocalEvent("Inserting value "+value, govec.GetDefaultLogOptions())

	/* Insert the Value */
	valueEntry := ValueEntry{value, logger.GetCurrentVC()}
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: valueEntry}}}}
	updateResult, err := db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("REPLICA "+strconv.Itoa(no+1)+": Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)

	util.PrintMsg(no, "Inserted value "+value)
}
