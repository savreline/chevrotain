package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Collection

// Post is
type Post struct {
	Title string `json:”title,omitempty”`
	Body  string `json:”body,omitempty”`
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database("chevrotain").Collection("kvs")
	InsertPost("1", "Hello")
	InsertPost("2", "Hello")
	defer client.Disconnect(ctx)
}

// InsertPost is https://www.mongodb.com/golang
func InsertPost(title string, body string) {
	post := Post{title, body}
	insertResult, err := db.InsertOne(context.TODO(), post)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
}
