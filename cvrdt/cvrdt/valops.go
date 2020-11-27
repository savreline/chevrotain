package main

import (
	"context"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertValue inserts value into the given key
func (t *RPCExt) InsertValue(args *util.RPCExtArgs, reply *int) error {
	InsertValueLocal(args.Key, args.Value, posCollection, nil)
	return nil
}

// RemoveValue removes value from the given key
func (t *RPCExt) RemoveValue(args *util.RPCExtArgs, reply *int) error {
	InsertValueLocal(args.Key, args.Value, negCollection, nil)
	return nil
}

// InsertValueLocal inserts the value in either positive collection (add) or negative collection (remove)
func InsertValueLocal(key string, value string, collection string, record *util.ValueEntry) {
	/* In no ready to go value entry record is supplied, tick the clock and make one */
	if record == nil {
		logger.LogLocalEvent("Inserting value "+value, govec.GetDefaultLogOptions())
		record = &util.ValueEntry{Value: value, Timestamp: logger.GetCurrentVC()}
	}

	/* Insert the Value */
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: record}}}}
	_, err := db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	util.PrintMsg(noStr, "Inserted value "+value)
}
