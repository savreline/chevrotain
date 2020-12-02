package main

import (
	"context"
	"fmt"

	"../util"
	"go.mongodb.org/mongo-driver/bson"
)

// Constants
const (
	collectionName = "kvs"
)

func insertKey(key string) {
	record := util.CmRecord{Name: key, Values: []string{}}
	_, err := db.Collection(collectionName).InsertOne(context.TODO(), record)
	if err != nil {
		fmt.Println(err)
	}
	if verbose == true {
		fmt.Println("Inserted Key " + key)
	}
}

func removeKey(key string) {
	filter := bson.D{{Key: "name", Value: key}}
	_, err := db.Collection(collectionName).DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}
	if verbose == true {
		fmt.Println("Removed Key " + key)
	}
}

func insertValue(key string, value string) {
	/* Define filters */
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: value}}}}

	/* Do the update */
	_, err := db.Collection(collectionName).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
	}
	if verbose == true {
		fmt.Println("Inserted Value " + value)
	}
}

func removeValue(key string, value string) {
	/* Define filters */
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$pull", Value: bson.D{
		{Key: "values", Value: value}}}}

	/* Do the delete */
	_, err := db.Collection(collectionName).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
	}
	if verbose == true {
		fmt.Println("Deleted Value " + value)
	}
}
