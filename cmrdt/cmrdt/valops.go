package main

import (
	"context"

	"../../util"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertValue inserts value into the given key
func (t *RPCExt) InsertValue(args *util.RPCExtArgs, reply *int) error {
	processExtCall(*args, IV)
	return nil
}

// RemoveValue givens the value from the given key
func (t *RPCExt) RemoveValue(args *util.RPCExtArgs, reply *int) error {
	processExtCall(*args, RV)
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
