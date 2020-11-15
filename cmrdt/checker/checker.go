package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"

	"../../util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Record is a DB Record
type Record struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func main() {
	drop := os.Args[1]
	_, _, dbPorts, err := util.ParseGroupMembersCVS("../driver/ports.csv", "")
	if err != nil {
		util.PrintErr(err)
	}
	noReplicas := len(dbPorts)
	cols := make([]*mongo.Collection, noReplicas)
	results := make([][]Record, noReplicas)

	/* Connect */
	for i, dbPort := range dbPorts {
		dbClient, _ := util.Connect(dbPort)
		cols[i] = dbClient.Database("chev").Collection("kvs")
		fmt.Println("Connected to DB on port " + dbPort)
	}

	// https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.Find
	// https://github.com/mongodb/mongo-go-driver
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	for i, col := range cols {
		cursor, err := col.Find(context.TODO(), bson.D{}, opts)
		if err != nil {
			log.Fatal(err)
		}
		if err = cursor.All(context.TODO(), &results[i]); err != nil {
			log.Fatal(err)
		}
		if drop == "1" {
			col.Drop(context.TODO())
		}
	}

	var result = true
	for i := 0; i < noReplicas-1; i++ {
		eqResult := testEq(results[i], results[i+1], i)
		fmt.Println("Comparison of", i+1, "to", i+2, "is", eqResult)
		if eqResult == false {
			result = false
		}
	}

	fmt.Println("Overall Result is", result)
}

// https://stackoverflow.com/questions/15311969/checking-the-equality-of-two-slices
func testEq(a, b []Record, no int) bool {
	var str1, str2 string
	result := true
	if (a == nil) != (b == nil) {
		result = false
	}
	if len(a) != len(b) {
		result = false
	}
	for i := range a {
		sort.Strings(a[i].Values)
		sort.Strings(b[i].Values)
		str1 = str1 + a[i].Name
		str2 = str2 + b[i].Name
		for _, val := range a[i].Values {
			str1 = str1 + "," + val
		}
		for _, val := range b[i].Values {
			str2 = str2 + "," + val
		}
		if !reflect.DeepEqual(a[i], b[i]) {
			result = false
		}
		str1 = str1 + "\n"
		str2 = str2 + "\n"
	}
	err := ioutil.WriteFile("Repl"+strconv.Itoa(no)+".csv", []byte(str1), 0644)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("Repl"+strconv.Itoa(no+1)+".csv", []byte(str2), 0644)
	if err != nil {
		panic(err)
	}
	return result
}
