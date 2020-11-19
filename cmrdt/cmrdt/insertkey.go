package main

/* In this file
1. InsertKey Ext RPC method
2. InsertKeyLocal (that works with the local db)
*/

import (
	"context"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *util.KeyArgs, reply *int) error {
	logger.StartBroadcast("OUT"+noStr+" InsKey "+args.Key, govec.GetDefaultLogOptions())
	InsertKeyLocal(args.Key)
	calls := broadcastInsert(args.Key, "")
	logger.StopBroadcast()
	waitForCallsToComplete(args.Key, "", calls)
	return nil
}

// InsertKeyLocal inserts the key into the local db
func InsertKeyLocal(key string) {
	filter := bson.D{{Key: "name", Value: "Keys"}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: key}}}}

	/* Update global keys entry */
	_, err := db.Collection("kvs").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(err)
	}
	// fmt.Printf("REPLICA "+strconv.Itoa(no)+": Matched %v documents and updated %v documents.\n",
	// 	updateResult.MatchedCount, updateResult.ModifiedCount)

	/* Insert entry for the given key */
	newRecord := Record{key, []string{}}
	_, err = db.Collection("kvs").InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}
	util.PrintMsg(no, "Inserted Key "+key)
}
