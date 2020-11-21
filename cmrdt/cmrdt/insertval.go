package main

import (
	"context"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"github.com/savreline/GoVector/govec/vclock"
	"go.mongodb.org/mongo-driver/bson"
)

// ValueArgs are the arguments to the InsertValue RPCInt call
type ValueArgs struct {
	Key, Value string
	Pid        string
	Timestamp  vclock.VClock
}

// InsertValue inserts value into the given key
func (t *RPCExt) InsertValue(args *ValueArgs, reply *int) error {
	logger.StartBroadcast("OUT"+noStr+" InsKey "+args.Key, govec.GetDefaultLogOptions())
	InsertValueLocal(args.Key, args.Value)
	calls := broadcastInsert(args.Key, args.Value, logger.GetCurrentVC())
	logger.StopBroadcast()
	waitForBroadcastToFinish(calls)
	return nil
}

// InsertValueLocal inserts the value into the local db
func InsertValueLocal(key string, value string) {
	/* Define filters */
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: value}}}}

	/* Do the update */
	_, err := db.Collection("kvs").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	util.PrintMsg(noStr, "Inserted Value "+value)
}
