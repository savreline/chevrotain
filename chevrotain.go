package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Collection

// KVPair is a key-value pair
type KVPair struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

func main() {
	/* Connect to MongoDB, reading port number from the command line
	As per https://www.mongodb.com/golang */
	port := os.Args[1]
	urlString := "mongodb://localhost:" + port + "/"

	client, err := mongo.NewClient(options.Client().ApplyURI(urlString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database("chev").Collection("kvs")

	InsertKey("1")
	InsertKey("2")
	InsertValue("1", "Hello")
	defer client.Disconnect(ctx)
}

// InsertKey inserts key with an empty array for values
func InsertKey(key string) {
	insertResult, err := db.InsertOne(context.TODO(), KVPair{key, []string{}})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted key with ID:", insertResult.InsertedID)
}

// InsertValue inserts value into the given key
func InsertValue(key string, value string) {
	filter := bson.D{{Key: key}}
	update := bson.D{
		{"$push", bson.D{
			{"value", value},
		}},
	}
	updateResult, err := db.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)
}
