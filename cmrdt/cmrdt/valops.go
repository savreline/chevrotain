package main

import (
	"context"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertValue inserts value into the given key
func (t *RPCExt) InsertValue(args *util.RPCExtArgs, reply *int) error {
	logger.StartBroadcast("OUT"+noStr+" InsVal "+args.Key+":"+args.Value, govec.GetDefaultLogOptions())
	InsertValueLocal(args.Key, args.Value)
	processExtCall(args, true)
	return nil
}

// RemoveValue remove the value frome the given key
func (t *RPCExt) RemoveValue(args *util.RPCExtArgs, reply *int) error {
	logger.StartBroadcast("OUT"+noStr+" RmvVal "+args.Key+":"+args.Value, govec.GetDefaultLogOptions())
	RemoveValueLocal(args.Key, args.Value)
	processExtCall(args, false)
	return nil
}

// InsertValueLocal inserts the value into the local db
func InsertValueLocal(key string, value string) {
	/* Define filters */
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: value}}}}

	/* Do the update */
	_, err := db.Collection(collectionName).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	if verbose == true {
		util.PrintMsg(noStr, "Inserted Value "+value)
	}
}

// RemoveValueLocal removes the value from the local db
func RemoveValueLocal(key string, value string) {
	/* Define filters */
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$pull", Value: bson.D{
		{Key: "values", Value: value}}}}

	/* Do the delete */
	_, err := db.Collection(collectionName).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	if verbose == true {
		util.PrintMsg(noStr, "Deleted Value "+value)
	}
}
