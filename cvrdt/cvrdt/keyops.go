package main

import (
	"context"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *util.RPCExtArgs, reply *int) error {
	InsertKeyLocal(args.Key, posCollection, nil)
	return nil
}

// RemoveKey removes the given key
func (t *RPCExt) RemoveKey(args *util.RPCExtArgs, reply *int) error {
	InsertKeyLocal(args.Key, negCollection, nil)
	return nil
}

// InsertKeyLocal inserts the key in either positive collection (add) or negative collection (remove)
func InsertKeyLocal(key string, collection string, record *util.CvRecord) {
	/* In no ready to go record is supplied, tick the clock and make one */
	if record == nil {
		logger.LogLocalEvent("Inserting Key "+key, govec.GetDefaultLogOptions())
		record = &util.CvRecord{Name: key, Timestamp: logger.GetCurrentVC(), Values: []util.ValueEntry{}}
	} else {
		if len(record.Values) == 0 {
			record.Values = []util.ValueEntry{}
		}
	}

	/* Update global keys entry */
	valEntry := util.ValueEntry{Value: key, Timestamp: record.Timestamp}
	filter := bson.D{{Key: "name", Value: "Keys"}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: valEntry}}}}
	_, err := db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Insert the Key */
	_, err = db.Collection(collection).InsertOne(context.TODO(), record)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Print to console */
	if collection == posCollection {
		util.PrintMsg(noStr, "Inserted Key "+key)
	} else {
		util.PrintMsg(noStr, "Removed Key "+key)
	}
}
