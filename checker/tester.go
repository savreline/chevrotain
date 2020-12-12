package main

import (
	"fmt"
	"os"
	"sort"

	"../util"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	/* Parse args and group membership info */
	drop := os.Args[1]
	impl := os.Args[2]
	ips, _, dbPorts, err := util.ParseGroupMembersCVS("../ports.csv", "")
	if err != nil {
		util.PrintErr("TESTER", err)
	}
	noReplicas := len(dbPorts)

	/* Init data structures */
	colsP := make([]*mongo.Collection, noReplicas)
	colsN := make([]*mongo.Collection, noReplicas)
	colsD := make([]*mongo.Collection, noReplicas)
	cols := make([]*mongo.Collection, noReplicas)
	results := make([][]util.SRecord, noReplicas)

	/* Connect to databases */
	for i, dbPort := range dbPorts {
		dbClient, _ := util.ConnectDb("TESTER", ips[i], dbPort)
		colsP[i] = dbClient.Database("chev").Collection("kvsp")
		colsN[i] = dbClient.Database("chev").Collection("kvsn")
		colsD[i] = dbClient.Database("chev").Collection("kvsd")
		cols[i] = dbClient.Database("chev").Collection("kvs")
		util.PrintMsg("TESTER", "Connected to DB on port "+dbPort)
	}

	/* Download results and save */
	if impl == "cv" {
		for i, col := range colsP {
			result := util.DownloadDState(col, "TESTER", drop)
			util.SaveDStateToCSV(result, i, "P")
		}
		for i, col := range colsN {
			result := util.DownloadDState(col, "TESTER", drop)
			util.SaveDStateToCSV(result, i, "N")
		}
	}
	if impl == "cm" {
		for i, col := range colsD {
			result := util.DownloadDState(col, "TESTER", drop)
			util.SaveDStateToCSV(result, i, "D")
		}
	}
	for i, col := range cols {
		results[i] = util.DownloadSState(col, "TESTER", drop)
		util.SaveSStateToCSV(results[i], i)
	}

	/* Check Equality */
	var i1, i2, i3 int
	for i := 0; i < noReplicas; i++ {
		if i != noReplicas-1 {
			i1 = i
			i2 = i + 1
			i3 = i + 2
		} else {
			i1 = 0
			i2 = noReplicas - 1
			i3 = noReplicas
		}
		errs, cnt := testEq(results[i1], results[i2])
		msg := fmt.Sprint("Number of diffs ", i1+1, " to ", i3, " is ", errs,
			" and percent is ", (cnt-errs)/cnt*100)
		util.PrintMsg("CHECKER", msg)
	}
}

// check databases for equality
// with advice from https://stackoverflow.com/questions/15311969/checking-the-equality-of-two-slices
func testEq(a, b []util.SRecord) (float32, float32) {
	var errs, cnt float32

	for i := range a {
		/* This signal a substatintial discprenacy, rough estimates of errors to finish off somehow */
		if i > len(b)-1 {
			errs += float32(len(a[i].Values))
			cnt += float32(len(a[i].Values))
			return errs, cnt
		}
		sort.Strings(a[i].Values)
		sort.Strings(b[i].Values)

		if a[i].Key != b[i].Key {
			errs++
		}
		cnt++

		for j, val := range a[i].Values {
			if j >= len(b[i].Values) || val != b[i].Values[j] {
				errs++
			}
			cnt++
		}
	}

	return errs, cnt
}
