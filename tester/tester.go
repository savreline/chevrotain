package main

import (
	"os"

	"../util"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	/* Parse args and group membership info */
	drop := os.Args[1]
	_, dbPorts, err := util.ParseGroupMembersCVS("../ports.csv", "")
	if err != nil {
		util.PrintErr("TESTER", err)
	}
	noReplicas := len(dbPorts)

	/* Init data structures */
	colsP := make([]*mongo.Collection, noReplicas)
	colsN := make([]*mongo.Collection, noReplicas)
	cols := make([]*mongo.Collection, noReplicas)

	/* Connect to databases */
	for i, dbPort := range dbPorts {
		dbClient, _ := util.ConnectDb("TESTER", dbPort)
		colsP[i] = dbClient.Database("chev").Collection("kvsp1")
		colsN[i] = dbClient.Database("chev").Collection("kvsn1")
		cols[i] = dbClient.Database("chev").Collection("kvs")
		util.PrintMsg("TESTER", "Connected to DB on port "+dbPort)
	}

	/* Download results and save */
	for i, col := range colsP {
		result := util.DownloadDState(col, "TESTER", drop)
		util.SaveDStateToCSV(result, i, "P")
	}
	for i, col := range colsN {
		result := util.DownloadDState(col, "TESTER", drop)
		util.SaveDStateToCSV(result, i, "N")
	}
	for i, col := range cols {
		result := util.DownloadSState(col, "TESTER", drop)
		util.SaveSStateToCSV(result, i)
	}
}
