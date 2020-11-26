package main

import (
	"context"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *util.RPCExtArgs, reply *int) error {
	InsertKeyLocal(args.Key, posCollection)
	return nil
}

// RemoveKey removes the given key
func (t *RPCExt) RemoveKey(args *util.RPCExtArgs, reply *int) error {
	InsertKeyLocal(args.Key, negCollection)
	return nil
}

// InsertKeyLocal inserts the key in either positive collection (add) or negative collection (remove)
func InsertKeyLocal(key string, collection string) {
	/* Tick the clock */
	logger.LogLocalEvent("Inserting Key "+key, govec.GetDefaultLogOptions())

	/* Update global keys entry */
	valueEntry := util.ValueEntry{Value: key, Timestamp: logger.GetCurrentVC()}
	filter := bson.D{{Key: "name", Value: "Keys"}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: valueEntry}}}}
	_, err := db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Insert the Key */
	newRecord := util.CvRecord{Name: key, Timestamp: logger.GetCurrentVC(), Values: []util.ValueEntry{}}
	_, err = db.Collection(collection).InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	util.PrintMsg(noStr, "Inserted Key "+key)
}
