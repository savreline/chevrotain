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
	InsertLocalRecord("Keys", args.Key, posCollection, nil)
	if delay > 0 {
		time.Sleep(time.Duration(util.GetRand(delay)) * time.Millisecond)
	}
	return nil
}

// RemoveKey removes the given key
func (t *RPCExt) RemoveKey(args *util.RPCExtArgs, reply *int) error {
	InsertLocalRecord("Keys", args.Key, negCollection, nil)
	if delay > 0 {
		time.Sleep(time.Duration(util.GetRand(delay)) * time.Millisecond)
	}
	return nil
}

// InsertValue inserts value into the given key
func (t *RPCExt) InsertValue(args *util.RPCExtArgs, reply *int) error {
	InsertLocalRecord(args.Key, args.Value, posCollection, nil)
	if delay > 0 {
		time.Sleep(time.Duration(util.GetRand(delay)) * time.Millisecond)
	}
	return nil
}

// RemoveValue removes value from the given key
func (t *RPCExt) RemoveValue(args *util.RPCExtArgs, reply *int) error {
	InsertLocalRecord(args.Key, args.Value, negCollection, nil)
	if delay > 0 {
		time.Sleep(time.Duration(util.GetRand(delay)) * time.Millisecond)
	}
	return nil
}

// InsertLocalRecord inserts the record in either positive collection (add) or negative collection (remove)
func InsertLocalRecord(key string, value string, collection string, record *util.ValueEntry) {
	/* In no ready to go record is supplied, tick the clock and make one */
	if record == nil {
		timestamp := logger.LogLocalEvent("Inserting key "+key+" value "+value, govec.GetDefaultLogOptions())
		record = &util.ValueEntry{Value: value, Timestamp: timestamp}
	}

	/* Look for the key (which could be "Keys") and push the record in */
	var filter = bson.D{{Key: "name", Value: key}}
	var update = bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: record}}}}

	/* If key entry is found, no need to insert the key entry as well */
	var res util.CvRecord
	err := db.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		keyEntry := &util.CvRecord{Name: key, Values: []util.ValueEntry{}}
		_, err := db.Collection(collection).InsertOne(context.TODO(), keyEntry)
		if err != nil {
			util.PrintErr(noStr, err)
		}
	}

	/* Do the main operation on the database */
	_, err = db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Add ticks */
	addTicks(record.Timestamp)

	/* Print to console */
	if verbose {
		if collection == posCollection && key == "Keys" {
			util.PrintMsg(noStr, "Inserted Key "+value)
		} else if collection == negCollection && key == "Keys" {
			util.PrintMsg(noStr, "Removed Key "+value)
		} else if collection == posCollection {
			util.PrintMsg(noStr, "Inserted Value "+value+" on key "+key)
		} else {
			util.PrintMsg(noStr, "Removed Value "+value+" on key "+key)
		}
	}
}
