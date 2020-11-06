package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"

	"../../util"
	"../cmrdt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	noReplicas, err := strconv.Atoi(os.Args[1])
	if err != nil {
		util.PrintErr(err)
	}
	dbPorts := [2]string{"27018", "27019"}
	ctx := make([]context.Context, noReplicas)
	cols := make([]*mongo.Collection, noReplicas)
	dbClients := make([]*mongo.Client, noReplicas)
	results := make([][]cmrdt.Record, noReplicas)

	/* Connect */
	for i, dbPort := range dbPorts {
		dbClients[i], ctx[i] = util.Connect(dbPort)
		cols[i] = dbClients[i].Database("chev").Collection("kvs")
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
		col.Drop(context.TODO())
	}

	var result = true
	for i := 0; i < noReplicas-1; i++ {
		eqResult := testEq(results[i], results[i+1])
		fmt.Println("Comparison of", i+1, "to", i+2, "is", eqResult)
		if eqResult == false {
			result = false
		}
	}

	fmt.Println("Overall Result is", result)
}

// https://stackoverflow.com/questions/15311969/checking-the-equality-of-two-slices
func testEq(a, b []cmrdt.Record) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		sort.Strings(a[i].Values)
		sort.Strings(b[i].Values)
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}
