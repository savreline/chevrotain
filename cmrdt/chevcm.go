package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"../util"

	"github.com/DistributedClocks/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database
var logger *govec.GoLog

// Record is a DB Record
type Record struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func main() {
	/* Connect to MongoDB, reading port number from the command line */
	port := os.Args[1]
	client, ctx := util.Connect(port)
	db = client.Database("chev")
	defer client.Disconnect(ctx)

	/* Initialize GoVector logger */
	logger = govec.InitGoVector("MyProcess", "LogFile", govec.GetDefaultConfig())

	/* Pre-allocate Keys entry */
	newRecord := Record{"Keys", []string{}}
	_, err := db.Collection("kvs").InsertOne(context.TODO(), newRecord)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Pre-allocated keys")

	InsertKey("1")
	InsertKey("2")
	InsertValue("1", "Hello")
	InsertValue("2", "Bye")
}

/****************************/
/*** 1: INSERT/REMOVE KEY ***/
/****************************/

// InsertKey inserts the given key with an empty array for values
func InsertKey(key string) {
	InsertKeyLocal(key)
}

// RemoveKey removes the given key
func RemoveKey(key string) {
	// TODO
}

// InsertKeyLocal inserts the key into the local db
func InsertKeyLocal(key string) {
	logger.LogLocalEvent("Inserting Key"+key, govec.GetDefaultLogOptions())
	filter := bson.D{{Key: "name", Value: "Keys"}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: key}}}}

	updateResult, err := db.Collection("kvs").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)

	newRecord := Record{key, []string{}}
	_, err = db.Collection("kvs").InsertOne(context.TODO(), newRecord)
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
	InsertValueLocal(key, value)
}

// RemoveValue removes value from the given key
func RemoveValue(key string, value string) {
	// TODO
}

// InsertValueLocal inserts the value into the local db
func InsertValueLocal(key string, value string) {
	logger.LogLocalEvent("Inserting value"+value, govec.GetDefaultLogOptions())
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: value}}}}

	updateResult, err := db.Collection("kvs").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)
}
