package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"../util"

	"github.com/DistributedClocks/GoVector/govec"
	"github.com/DistributedClocks/GoVector/govec/vclock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	posCollection = "kvsp"
	negCollection = "kvsn"
)

var db *mongo.Database
var logger *govec.GoLog

// Record is a DB Record
type Record struct {
	Name   string        `json:"name"`
	Time   vclock.VClock `json:"time"`
	Values []ValueEntry  `json:"values"`
}

// ValueEntry is a value along with the timestamp
type ValueEntry struct {
	Value string        `json:"name"`
	Time  vclock.VClock `json:"time"`
}

func main() {
	/* Connect to MongoDB, reading port number from the command line */
	port := os.Args[1]
	client, ctx := util.Connect(port)
	db = client.Database("chev")
	defer client.Disconnect(ctx)

	/* Initialize GoVector logger */
	logger = govec.InitGoVector("Local", "LogFile", govec.GetDefaultConfig())

	/* Pre-allocate global keys entries */
	newRecord := Record{"Keys", logger.GetCurrentVC(), []ValueEntry{}}
	_, err := db.Collection(posCollection).InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println("Inserted global keys record into positive table")

	_, err = db.Collection(negCollection).InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println("Inserted global keys record into negative table")

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
	/* Update global keys entry */
	logger.LogLocalEvent("Inserting Key"+key, govec.GetDefaultLogOptions())
	valueEntry := ValueEntry{key, logger.GetCurrentVC()}
	filter := bson.D{{Key: "name", Value: "Keys"}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: valueEntry}}}}

	updateResult, err := db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)

	/* Setup a record for the current key */
	newRecord := Record{key, logger.GetCurrentVC(), []ValueEntry{}}
	_, err = db.Collection(collection).InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
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
	logger.LogLocalEvent("Inserting value"+value, govec.GetDefaultLogOptions())
	valueEntry := ValueEntry{value, logger.GetCurrentVC()}
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: valueEntry}}}}

	updateResult, err := db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)
}
