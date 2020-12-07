package main

import (
	"context"

	"../util"
	"go.mongodb.org/mongo-driver/bson"
)

// insert the key into the local db
func insertKey(key string) {
	record := util.SRecord{Key: key, Values: []string{}}
	_, err := db.Collection(sCollection).InsertOne(context.TODO(), record)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	if verbose {
		util.PrintMsg(noStr, "Inserted Key "+key)
	}
}

// insert the given value into the local db
func insertValue(key string, value string) {
	/* Define filters */
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: value}}}}

	/* Do the update */
	_, err := db.Collection(sCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	if verbose {
		util.PrintMsg(noStr, "Inserted Value "+value)
	}
}

// removes the given key from the local db
func removeKey(key string) {
	filter := bson.D{{Key: "name", Value: key}}
	_, err := db.Collection(sCollection).DeleteOne(context.TODO(), filter)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	if verbose {
		util.PrintMsg(noStr, "Deleted Key "+key)
	}
}

// removes the given value from the local db
func removeValue(key string, value string) {
	/* Define filters */
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$pull", Value: bson.D{
		{Key: "values", Value: value}}}}

	/* Do the delete */
	_, err := db.Collection(sCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	if verbose {
		util.PrintMsg(noStr, "Deleted Value "+value)
	}
}
