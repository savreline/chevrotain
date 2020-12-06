package main

import (
	"context"

	"../util"
	"go.mongodb.org/mongo-driver/bson"
)

// inserts the key or value into the local database along with the provided unique id
func insert(key string, value string, id int) {

	/* Check if the document to be updated exists, if not, make one */
	var dbResult util.DRecord
	filter := bson.D{{Key: "key", Value: key}}
	err := db.Collection(dCollection).FindOne(context.TODO(), filter).Decode(&dbResult)
	if err != nil {
		doc := util.DDoc{Key: key, Values: []util.DRecord{}}
		_, err = db.Collection(dCollection).InsertOne(context.TODO(), doc)
		if err != nil {
			util.PrintErr(noStr, err)
		}
	}

	/* Update the document */
	record := util.DRecord{Value: value, ID: id}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: record}}}}
	_, err = db.Collection(dCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Print to console */
	if verbose && key == "Keys" {
		util.PrintMsg(noStr, "Inserted Key "+value)
	} else {
		util.PrintMsg(noStr, "Inserted Value "+value+" on key "+key)
	}
}

// removes all instances of keys and values with the given ids
func remove(key string, value string, ids []int) {
	/* Do the remove */
	filter := bson.D{{Key: "key", Value: key}}
	update := bson.D{{Key: "$pull", Value: bson.D{
		{Key: "values", Value: bson.D{
			{Key: "id", Value: bson.D{
				{Key: "$in", Value: ids}}}}}}}}
	_, err := db.Collection(dCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Print to console */
	if verbose && key == "Keys" {
		util.PrintMsg(noStr, "Removed Key "+value)
	} else {
		util.PrintMsg(noStr, "Removed Value "+value+" on key "+key)
	}
}

// computes the set of all currently existing ids for the given key and value
func computeRemovalSet(key string, value string) []int {
	/* Find the matching document */
	var dbResult util.DDoc
	filter := bson.D{{Key: "key", Value: key}}
	err := db.Collection(dCollection).FindOne(context.TODO(), filter).Decode(&dbResult)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Extract ids from the mathcing document */
	var results []int
	for _, record := range dbResult.Values {
		if record.Value == value {
			results = append(results, record.ID)
		}
	}
	return results
}

// generates the "lookup" view collection of the database
func lookup() {
	/* Download state */
	state := util.DownloadDState(db.Collection(dCollection), "TESTER", "0")

	/* Download the "keys" document */
	var keysDoc util.DDoc
	filter := bson.D{{Key: "key", Value: "Keys"}}
	err := db.Collection(dCollection).FindOne(context.TODO(), filter).Decode(&keysDoc)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Insert keys */
	for _, record := range keysDoc.Values {
		util.InsertSKey(db.Collection(sCollection), noStr, record.Value)
	}

	/* Insert values, if the corresponding key exists */
	for _, doc := range state {
		for _, record := range doc.Values {
			if util.CheckMembership(keysDoc.Values, doc.Key) {
				util.InsertSValue(db.Collection(sCollection), noStr, doc.Key, record.Value)
			}
		}
	}
}
