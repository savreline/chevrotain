package util

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SRecord is a static Db record
type SRecord struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// DDoc is a dynamic Db document
type DDoc struct {
	Key    string    `json:"key"`
	Values []DRecord `json:"values"`
}

// DRecord is a dynamic Db record (id is timestamp in the case of CvRDT)
type DRecord struct {
	Value string `json:"value"`
	ID    int    `json:"id"`
}

// DownloadDState downloads the contents of any dynamic collection
func DownloadDState(col *mongo.Collection, who string, drop string) []DDoc {
	var res []DDoc

	/* Download all key docs */
	opts := options.Find().SetSort(bson.D{{Key: "key", Value: 1}})
	cursor, err := col.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		PrintErr(who, err)
	}

	/* Save downloaded info into res */
	if err = cursor.All(context.TODO(), &res); err != nil {
		PrintErr(who, err)
	}

	/* Drop the collection if asked */
	if drop == "1" {
		col.Drop(context.TODO())
	}
	return res
}

// DownloadSState downloads the contents of any static collection
func DownloadSState(col *mongo.Collection, who string, drop string) []SRecord {
	var result []SRecord

	opts := options.Find().SetSort(bson.D{{Key: "key", Value: 1}})
	cursor, err := col.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		PrintErr(who, err)
	}
	if err = cursor.All(context.TODO(), &result); err != nil {
		PrintErr(who, err)
	}
	if drop == "1" {
		col.Drop(context.TODO())
	}
	return result
}

// PrintDState prints a dynamic state to the console
func PrintDState(state []DDoc) {
	for _, doc := range state {
		fmt.Println(doc)
	}
	fmt.Println()
}

// PrintSState prints a static state to the console
func PrintSState(state []SRecord) {
	for _, record := range state {
		fmt.Println(record)
	}
	fmt.Println()
}

// InsertSKey inserts the given key into the static collection
func InsertSKey(col *mongo.Collection, who string, key string) {
	/* Check if the record exists */
	var dbResult SRecord
	filter := bson.D{{Key: "key", Value: key}}
	err := col.FindOne(context.TODO(), filter).Decode(&dbResult)

	/* Do the insert */
	if err != nil {
		record := SRecord{Key: key, Values: []string{}}
		_, err := col.InsertOne(context.TODO(), record)
		if err != nil {
			PrintErr(who, err)
		}
	}
}

// InsertSValue inserts the given value into the static collection
func InsertSValue(col *mongo.Collection, who string, key string, value string) {
	/* Check if the record exists */
	var dbResult SRecord
	filter := bson.D{{Key: "key", Value: key},
		{Key: "values", Value: value}}
	err := col.FindOne(context.TODO(), filter).Decode(&dbResult)

	if err != nil { // error exists, so didn't find it, so insert
		/* Check if the document to be updated exists, if not, make one */
		filter = bson.D{{Key: "key", Value: key}}
		err = col.FindOne(context.TODO(), filter).Decode(&dbResult)
		if err != nil {
			keyEntry := &SRecord{Key: key, Values: []string{}}
			_, err := col.InsertOne(context.TODO(), keyEntry)
			if err != nil {
				PrintErr(who, err)
			}
		}

		/* Do the update */
		filter := bson.D{{Key: "key", Value: key}}
		update := bson.D{{Key: "$push", Value: bson.D{
			{Key: "values", Value: value}}}}
		_, err := col.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			PrintErr(who, err)
		}
	}
}

// CheckMembership return true if an entry with the given value is found in a slice of records
func CheckMembership(arr []DRecord, value string) bool {
	for _, record := range arr {
		if record.Value == value {
			return true
		}
	}
	return false
}
