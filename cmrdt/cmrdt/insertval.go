package main

import (
	"context"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertValue inserts value into the given key
func (t *RPCExt) InsertValue(args *util.ValueArgs, reply *int) error {
	logger.StartBroadcast("OUT"+noStr+" InsKey "+args.Key, govec.GetDefaultLogOptions())
	opNode := OpNode{
		Type:      IV,
		Key:       args.Key,
		Value:     args.Value,
		Timestamp: logger.GetCurrentVC().Copy(),
		Pid:       noStr,
		ConcOp:    false}
	logger.StopBroadcast()
	addToQueue(opNode)
	calls := broadcast(opNode)
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
