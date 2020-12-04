package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

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
		colsP[i] = dbClient.Database("chev").Collection("kvsp")
		colsN[i] = dbClient.Database("chev").Collection("kvsn")
		cols[i] = dbClient.Database("chev").Collection("kvs")
		util.PrintMsg("TESTER", "Connected to DB on port "+dbPort)
	}

	/* Download results and save */
	for i, col := range colsP {
		result := util.DownloadCvState(col, "TESTER", drop)
		saveCvToCSV(result, i, "P")
	}
	for i, col := range colsN {
		result := util.DownloadCvState(col, "TESTER", drop)
		saveCvToCSV(result, i, "N")
	}
}

// save Cv data to CSV
func saveCvToCSV(a []util.CvDoc, no int, pn string) {
	var str string

	for i := range a {
		// sort.Strings(a[i].Values)
		str = str + a[i].Key
		for _, val := range a[i].Values {
			str = str + "," + fmt.Sprint(val)
		}
		str = str + "\n"
	}

	/* Write to CSV */
	err := ioutil.WriteFile("Repl"+pn+strconv.Itoa(no)+".csv", []byte(str), 0644)
	if err != nil {
		util.PrintErr("CHECKER", err)
	}
}
