package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"../../util"
	"go.mongodb.org/mongo-driver/mongo"
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

	/* Connect */
	for i, dbPort := range dbPorts {
		dbClient, _ := util.Connect("CHECKER", dbPort)
		colsP[i] = dbClient.Database("chev").Collection("kvsp")
		colsN[i] = dbClient.Database("chev").Collection("kvsn")
		util.PrintMsg("CHECKER", "Connected to DB on port "+dbPort)
	}

	/* Download results and save */
	for i, col := range colsP {
		result := util.DownloadCvState(col, drop)
		saveToCSV(result, i, "P")
	}
	for i, col := range colsN {
		result := util.DownloadCvState(col, drop)
		saveToCSV(result, i, "N")
	}
}

// save to CSV
func saveToCSV(a []util.CvRecord, no int, pn string) {
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
	err := ioutil.WriteFile("Repl"+pn+strconv.Itoa(no)+".csv", []byte(str), 0644)
	if err != nil {
		util.PrintErr("CHECKER", err)
	}
}
