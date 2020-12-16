package util

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"

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
func DownloadDState(db *mongo.Database, who string, name string, drop string) []DDoc {
	var res []DDoc

	/* Download all key docs */
	opts := options.Find().SetSort(bson.D{{Key: "key", Value: 1}})
	cursor, err := db.Collection(name).Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		PrintErr(who, err)
	}

	/* Save downloaded info into res */
	if err = cursor.All(context.TODO(), &res); err != nil {
		PrintErr(who, err)
	}

	/* Drop the collection if asked and initialize a new one */
	if drop == "1" {
		db.Collection(name).Drop(context.TODO())
		CreateCollection(db, who, name)
	}
	return res
}

// DownloadSState downloads the contents of any static collection
func DownloadSState(db *mongo.Database, who string, drop string) []SRecord {
	var result []SRecord
	name := "kvs"

	/* Download all key docs */
	opts := options.Find().SetSort(bson.D{{Key: "key", Value: 1}})
	cursor, err := db.Collection(name).Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		PrintErr(who, err)
	}

	/* Save downloaded info into res */
	if err = cursor.All(context.TODO(), &result); err != nil {
		PrintErr(who, err)
	}

	/* Drop the collection if asked and initialize a new one */
	if drop == "1" {
		db.Collection(name).Drop(context.TODO())
		CreateCollection(db, who, name)
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

// SaveDStateToCSV saves dynamic state to CSV
func SaveDStateToCSV(state []DDoc, no int, pn string) {
	var str string

	for _, doc := range state {
		// sort.Strings(a[i].Values)
		str = str + doc.Key
		for _, val := range doc.Values {
			str = str + "," + fmt.Sprint(val)
		}
		str = str + "\n"
	}

	/* Write to CSV */
	err := ioutil.WriteFile("Repl"+pn+strconv.Itoa(no)+".csv", []byte(str), 0644)
	if err != nil {
		PrintErr("CHECKER", err)
	}
}

// SaveSStateToCSV saves static state to CSV
func SaveSStateToCSV(state []SRecord, no int) {
	var str string

	for _, record := range state {
		sort.Strings(record.Values)
		str = str + record.Key
		for _, val := range record.Values {
			str = str + "," + fmt.Sprint(val)
		}
		str = str + "\n"
	}

	/* Write to CSV */
	err := ioutil.WriteFile("Repl"+strconv.Itoa(no)+".csv", []byte(str), 0644)
	if err != nil {
		PrintErr("CHECKER", err)
	}
}

// DownloadMainTestRef returns the data contained in the CSV reference file as an
// array of SRecords https://stackoverflow.com/questions/24999079/reading-csv-file-in-go
func DownloadMainTestRef() []SRecord {
	f, err := os.Open("MainTestRef.csv")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	res := []SRecord{}

	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				return res
			}
		}

		var record = SRecord{Key: row[0], Values: []string{}}
		for i, element := range row {
			if i != 0 {
				record.Values = append(record.Values, element)
			}
		}
		res = append(res, record)
	}
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

// CreateCollection checks if the specified collection exists, and creates one if that is not
// the case https://stackoverflow.com/questions/46293070/how-to-check-if-collection-exists-or-not-mongodb-golang
func CreateCollection(db *mongo.Database, who string, name string) {
	cols, err := db.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		PrintErr(who, err)
	}

	found := false
	for _, col := range cols {
		if col == name {
			found = true
			break
		}
	}

	if !found {
		err = db.CreateCollection(context.TODO(), name)
		if err != nil {
			PrintErr(who, err)
		}
	}
}
