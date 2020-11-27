package main

import (
	"context"

	"../../util"
)

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *util.RPCExtArgs, reply *int) error {
	processExtCall(*args, IK)
	return nil
}

// RemoveKey removes the given key
func (t *RPCExt) RemoveKey(args *util.RPCExtArgs, reply *int) error {
	processExtCall(*args, RK)
	return nil
}

// InsertKeyLocal inserts the key into the local db
func InsertKeyLocal(key string) {
	record := util.CmRecord{Name: key, Values: []string{}}
	_, err := db.Collection("kvs").InsertOne(context.TODO(), record)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	util.PrintMsg(noStr, "Inserted Key "+key)
}
