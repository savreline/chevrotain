package main

import (
	"context"

	"../../util"
	"github.com/savreline/GoVector/govec"
)

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *util.KeyArgs, reply *int) error {
	logger.StartBroadcast("OUT"+noStr+" InsKey "+args.Key, govec.GetDefaultLogOptions())
	InsertKeyLocal(args.Key)
	opNode := OpNode{
		Type:      IK,
		Key:       args.Key,
		Value:     "",
		Timestamp: logger.GetCurrentVC(),
		Pid:       noStr,
		ConcOp:    false}
	calls := broadcast(opNode)
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
