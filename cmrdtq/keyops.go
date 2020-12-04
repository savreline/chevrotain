package main

import (
	"context"

	"../../util"
	"go.mongodb.org/mongo-driver/bson"
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
