package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database
var posCollection = "kvsp"
var negCollection = "kvsn"

// Record is a DB Record
type Record struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func main() {
	/* Connect to MongoDB, reading port number from the command line */
	port := os.Args[1]
	client, ctx := connect(port)
	db = client.Database("chev")
	defer client.Disconnect(ctx)

	/* Pre-allocate Keys entry */
	_, err := db.Collection(posCollection).InsertOne(context.TODO(), Record{"Keys", []string{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted keys record into positive table")

	_, err = db.Collection(negCollection).InsertOne(context.TODO(), Record{"Keys", []string{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted keys record into negative table")

	InsertKey("1")
	InsertKey("2")
	InsertValue("1", "Hello")
	InsertValue("2", "Bye")
}

/******************************************/
/*** 1: INSERT/REMOVE KEY LOCAL METHODS ***/
/******************************************/

// InsertKey inserts the given key with an empty array for values
func InsertKey(key string) {
	InsertKeyHelper(key, posCollection)
}

// RemoveKey removes the given key
func RemoveKey(key string) {
	InsertKeyHelper(key, negCollection)
}

// InsertKeyHelper inserts the key in either positive collection (add) or negative collection (remove)
func InsertKeyHelper(key string, collection string) {
	filter := bson.D{{Key: "name", Value: "Keys"}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: key}}}}

	updateResult, err := db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)

	_, err = db.Collection(collection).InsertOne(context.TODO(), Record{key, []string{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted key", key)
}

/********************************************/
/*** 2: INSERT/REMOVE VALUE LOCAL METHODS ***/
/********************************************/

// InsertValue inserts value into the given key
func InsertValue(key string, value string) {
	InsertValueHelper(key, value, posCollection)
}

// RemoveValue removes value from the given key
func RemoveValue(key string, value string) {
	InsertValueHelper(key, value, negCollection)
}

// InsertValueHelper inserts the value in either positive collection (add) or negative collection (remove)
func InsertValueHelper(key string, value string, collection string) {
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: value}}}}

	updateResult, err := db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)
}
