package main

import (
	"context"
	"time"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *util.RPCExtArgs, reply *int) error {
	logger.StartBroadcast("OUT"+noStr+" InsKey "+args.Key, govec.GetDefaultLogOptions())
	InsertKeyLocal(args.Key)
	processExtCall(args, true)
	return nil
}

// RemoveKey removes the given key
func (t *RPCExt) RemoveKey(args *util.RPCExtArgs, reply *int) error {
	logger.StartBroadcast("OUT"+noStr+" RmvKey "+args.Key, govec.GetDefaultLogOptions())
	RemoveKeyLocal(args.Key)
	processExtCall(args, false)
	return nil
}

func processExtCall(args *util.RPCExtArgs, insert bool) {
	calls := broadcast(args, insert)
	logger.StopBroadcast()
	waitForCallsToComplete(calls)
	if delay > 0 {
		time.Sleep(time.Duration(util.GetRand(delay)) * time.Millisecond)
	}
}

// InsertKeyLocal inserts the key into the local db
func InsertKeyLocal(key string) {
	record := util.CmRecord{Name: key, Values: []string{}}
	_, err := db.Collection(collectionName).InsertOne(context.TODO(), record)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	if verbose == true {
		util.PrintMsg(noStr, "Inserted Key "+key)
	}
}

// RemoveKeyLocal removes the key from the local db
func RemoveKeyLocal(key string) {
	filter := bson.D{{Key: "name", Value: key}}
	_, err := db.Collection(collectionName).DeleteOne(context.TODO(), filter)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	if verbose == true {
		util.PrintMsg(noStr, "Deleted Key "+key)
	}
}
