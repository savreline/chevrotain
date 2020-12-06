package main

import (
	"context"

	"../util"
	"go.mongodb.org/mongo-driver/bson"
)

// inserts the record in either positive collection (add) or negative collection (remove)
func insertLocalRecord(key string, value string, collection string, record *util.DRecord) {
	/* In no ready to go record is supplied, tick the clock and make one,
	otherwise check if an exact identical entry already exists */
	if record == nil {
		record = &util.DRecord{Value: value, ID: clock}
	} else {
		var res util.DDoc
		filter := bson.D{{Key: "key", Value: key},
			{Key: "values", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "value", Value: record.Value},
					{Key: "timestamp", Value: record.ID}}}}}}
		err := db.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)
		if err == nil { // found the record, no need to do anything
			return
		}
	}

	/* Tick the clock */
	clock++

	/* Check if the document to be updated exists, if not, make one */
	var dbResult util.DRecord
	filter := bson.D{{Key: "key", Value: key}}
	err := db.Collection(collection).FindOne(context.TODO(), filter).Decode(&dbResult)
	if err != nil {
		keyEntry := &util.DDoc{Key: key, Values: []util.DRecord{}}
		_, err := db.Collection(collection).InsertOne(context.TODO(), keyEntry)
		if err != nil {
			util.PrintErr(noStr, err)
		}
	}

	/* Update the document */
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: record}}}}
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

// insert key into the static collection
func insertKey(key string) {
	util.InsertSKey(db.Collection(sCollection), noStr, key)
}

// insert value into the static collection
func insertValue(key string, value string) {
	util.InsertSValue(db.Collection(sCollection), noStr, key, value)
}

// removes key from the static collection
func removeKey(key string) {
	filter := bson.D{{Key: "key", Value: key}}
	_, err := db.Collection(sCollection).DeleteOne(context.TODO(), filter)
	if err != nil {
		util.PrintErr(noStr, err)
	}
}

// removes value from the static collection
func removeValue(value string, key string) {
	filter := bson.D{{Key: "key", Value: key}}
	update := bson.D{{Key: "$pull", Value: bson.D{
		{Key: "values", Value: value}}}}
	_, err := db.Collection(sCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}
}
