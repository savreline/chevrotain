package main

/* In this file
1. InsertKey Ext RPC method
2. InsertKeyLocal (that works with the local db)
3. InsertKeyGlobal (that broadcats the event) methods
*/

import (
	"context"
	"fmt"
	"strconv"

	"../../util"
	"github.com/DistributedClocks/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *util.KeyArgs, reply *int) error {
	InsertKeyLocal(args.Key)
	InsertKeyGlobal(args.Key)
	return nil
}

// InsertKeyLocal inserts the key into the local db
func InsertKeyLocal(key string) {
	logger.LogLocalEvent("Inserting Key"+key, govec.GetDefaultLogOptions())
	filter := bson.D{{Key: "name", Value: "Keys"}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: key}}}}

	/* Update global keys entry */
	updateResult, err := db.Collection("kvs").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Printf("REPLICA "+strconv.Itoa(no)+": Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)

	/* Insert entry for the given key */
	newRecord := Record{key, []string{}}
	_, err = db.Collection("kvs").InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}
	util.PrintMsg(no, "Inserted Key "+key)
}

// InsertKeyGlobal broadcasts the insertKey operation to other replicas
func InsertKeyGlobal(key string) {
	var result int
	var destNo int
	var flag = false

	for i, client := range conns {
		if i == no {
			flag = true
		}

		if client != nil {
			if flag {
				destNo = i + 1
			} else {
				destNo = i
			}
			fmt.Println("Sending RPC", no, "->", destNo)
			err := client.Call("RPCInt.InsertKeyRPC", util.KeyArgs{Key: key}, &result)
			if err != nil {
				util.PrintErr(err)
			}
		}
	}
}
