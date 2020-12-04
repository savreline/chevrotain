package main

import (
	"context"
	"time"

	"../util"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *util.RPCExtArgs, reply *int) error {
	InsertLocalRecord("Keys", args.Key, posCollection, nil)
	emulateDelay()
	return nil
}

// RemoveKey removes the given key
func (t *RPCExt) RemoveKey(args *util.RPCExtArgs, reply *int) error {
	InsertLocalRecord("Keys", args.Key, negCollection, nil)
	emulateDelay()
	return nil
}

// InsertValue inserts value into the given key
func (t *RPCExt) InsertValue(args *util.RPCExtArgs, reply *int) error {
	InsertLocalRecord(args.Key, args.Value, posCollection, nil)
	emulateDelay()
	return nil
}

// RemoveValue removes value from the given key
func (t *RPCExt) RemoveValue(args *util.RPCExtArgs, reply *int) error {
	InsertLocalRecord(args.Key, args.Value, negCollection, nil)
	emulateDelay()
	return nil
}

// emulates link delay in all RPC responses
func emulateDelay() {
	if delay > 0 {
		time.Sleep(time.Duration(util.GetRand(delay)) * time.Millisecond)
	}
}

// InsertLocalRecord inserts the record in either positive collection (add) or negative collection (remove)
func InsertLocalRecord(key string, value string, collection string, record *util.CvRecord) {
	/* In no ready to go record is supplied, tick the clock and make one,
	otherwise check if an exact identical entry already exists */
	if record == nil {
		record = &util.CvRecord{Value: value, Timestamp: clock}
	} else {
		var res util.CvDoc
		filter := bson.D{{Key: "key", Value: key},
			{Key: "values", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "value", Value: record.Value},
					{Key: "timestamp", Value: record.Timestamp}}}}}}
		err := db.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)
		if err == nil { // found the record, no need to do anything
			return
		}
	}

	/* Tick the clock */
	clock++

	/* Look for the key document (which could be "Keys") and push the record in */
	var filter = bson.D{{Key: "key", Value: key}}
	var update = bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: record}}}}

	/* If key document is found, no need to insert the document as well */
	var res util.CvRecord
	err := db.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		keyEntry := &util.CvDoc{Key: key, Values: []util.CvRecord{}}
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
