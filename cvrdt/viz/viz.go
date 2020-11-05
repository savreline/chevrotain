package main

import (
	"context"
	"fmt"
	"log"

	"../../util"
	"../cvrdt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	_, dbPorts, err := util.ParseGroupMembersCVS("../driver/ports.csv", "")
	if err != nil {
		util.PrintErr(err)
	}
	noReplicas := len(dbPorts)
	ctx := make([]context.Context, noReplicas)
	colsP := make([]*mongo.Collection, noReplicas)
	colsN := make([]*mongo.Collection, noReplicas)
	dbClients := make([]*mongo.Client, noReplicas)
	resultsP := make([][]cvrdt.Record, noReplicas)
	resultsN := make([][]cvrdt.Record, noReplicas)

	/* Connect */
	for i, dbPort := range dbPorts {
		dbClients[i], ctx[i] = util.Connect(dbPort)
		colsP[i] = dbClients[i].Database("chev").Collection("kvsp")
		colsN[i] = dbClients[i].Database("chev").Collection("kvsn")
		fmt.Println("Connected to DB on port " + dbPort)
	}

	/* Download Data Into Results Slice */
	getData(colsP, resultsP)
	getData(colsN, resultsN)

	/* Show Data */
	fmt.Println()
	fmt.Println("*** POSITIVE SET DATA ***")
	showData(resultsP)
	fmt.Println("*** NEGATIVE SET DATA ***")
	showData(resultsN)

	/* Clear Database? */
	// https://gist.github.com/albrow/5882501
	fmt.Println("Clear Databases? (Press 1 for yes)")
	var response string
	_, err = fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	if response == "1" {
		for i := 0; i < noReplicas; i++ {
			colsP[i].Drop(context.TODO())
			colsN[i].Drop(context.TODO())
		}
	}
}

// https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.Find
// https://github.com/mongodb/mongo-go-driver
func getData(cols []*mongo.Collection, results [][]cvrdt.Record) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	for i, col := range cols {
		cursor, err := col.Find(context.TODO(), bson.D{}, opts)
		if err != nil {
			log.Fatal(err)
		}
		if err = cursor.All(context.TODO(), &results[i]); err != nil {
			log.Fatal(err)
		}
	}
}

func showData(results [][]cvrdt.Record) {
	for i, resultSlice := range results {
		fmt.Println("REPLICA", i+1)
		for _, result := range resultSlice {
			fmt.Println(result)
		}
		fmt.Println()
	}
}
