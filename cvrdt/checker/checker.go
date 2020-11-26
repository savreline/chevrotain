package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"../../util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	/* Parse args, initialize data structures */
	drop := os.Args[1]
	_, dbPorts, err := util.ParseGroupMembersCVS("../driver/ports.csv", "")
	if err != nil {
		util.PrintErr("CHECKER", err)
	}
	noReplicas := len(dbPorts)
	colsP := make([]*mongo.Collection, noReplicas)
	colsN := make([]*mongo.Collection, noReplicas)
	resultsP := make([][]util.CvRecord, noReplicas)
	resultsN := make([][]util.CvRecord, noReplicas)

	/* Connect */
	for i, dbPort := range dbPorts {
		dbClient, _ := util.Connect("CHECKER", dbPort)
		colsP[i] = dbClient.Database("chev").Collection("kvsp")
		colsN[i] = dbClient.Database("chev").Collection("kvsn")
		util.PrintMsg("CHECKER", "Connected to DB on port "+dbPort)
	}

	/* Download results and save */
	downloadResults(colsP, resultsP, drop)
	for i, result := range resultsP {
		saveToCSV(result, i)
	}
	downloadResults(colsN, resultsN, drop)
	// for i, result := range resultsN {
	// 	saveToCSV(result, i)
	// }
}

// https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.Find
// https://github.com/mongodb/mongo-go-driver
func downloadResults(cols []*mongo.Collection, results [][]util.CvRecord, drop string) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	for i, col := range cols {
		cursor, err := col.Find(context.TODO(), bson.D{}, opts)
		if err != nil {
			util.PrintErr("CHECKER", err)
		}
		if err = cursor.All(context.TODO(), &results[i]); err != nil {
			util.PrintErr("CHECKER", err)
		}
		if drop == "1" {
			col.Drop(context.TODO())
		}
	}
}

// save to CSV
func saveToCSV(a []util.CvRecord, no int) {
	var str string

	for i := range a {
		// sort.Strings(a[i].Values)
		str = str + a[i].Name

		for _, val := range a[i].Values {
			str = str + "," + val.Value
		}

		str = str + "\n"
		str = str + fmt.Sprint(a[i].Timestamp)

		for _, val := range a[i].Values {
			str = str + "," + fmt.Sprint(val.Timestamp)
		}

		str = str + "\n"
	}

	/* Write to CSV */
	err := ioutil.WriteFile("Repl"+strconv.Itoa(no)+".csv", []byte(str), 0644)
	if err != nil {
		util.PrintErr("CHECKER", err)
	}
}
