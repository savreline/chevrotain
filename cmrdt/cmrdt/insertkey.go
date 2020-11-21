package main

import (
	"context"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"github.com/savreline/GoVector/govec/vclock"
)

// KeyArgs are the arguments to the InsertKey RPCInt call
type KeyArgs struct {
	Key       string
	Pid       string
	Timestamp vclock.VClock
}

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *KeyArgs, reply *int) error {
	logger.StartBroadcast("OUT"+noStr+" InsKey "+args.Key, govec.GetDefaultLogOptions())
	InsertKeyLocal(args.Key)
	calls := broadcastInsert(args.Key, "", logger.GetCurrentVC())
	logger.StopBroadcast()
	waitForBroadcastToFinish(calls)
	return nil
}

// InsertKeyLocal inserts the key into the local db
func InsertKeyLocal(key string) {
	newRecord := util.Record{Name: key, Values: []string{}}
	_, err := db.Collection("kvs").InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	util.PrintMsg(noStr, "Inserted Key "+key)
}
