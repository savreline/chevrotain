package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
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
	cols := make([]*mongo.Collection, noReplicas)
	results := make([][]util.CmRecord, noReplicas)

	/* Connect */
	for i, dbPort := range dbPorts {
		dbClient, _ := util.Connect("CHECKER", dbPort)
		cols[i] = dbClient.Database("chev").Collection("kvs")
		util.PrintMsg("CHECKER", "Connected to DB on port "+dbPort)
	}

	/* Get Snapshot of Database's State */
	for i, col := range cols {
		results[i] = util.DownloadCmState(col, drop)
	}

	/* Check Equality */
	for i := 0; i < noReplicas-1; i++ {
		errs, cnt := testEq(results[i], results[i+1], i)
		msg := fmt.Sprint("Number of diffs ", i+1, " to ", i+2, " is ", errs, " and percent is ", (cnt-errs)/cnt*100)
		util.PrintMsg("CHECKER", msg)
	}
}

// with advice from https://stackoverflow.com/questions/15311969/checking-the-equality-of-two-slices
func testEq(a, b []util.CmRecord, no int) (float32, float32) {
	var str1, str2 string
	var errs, cnt float32

	/* Check Equality */
	for i := range a {
		sort.Strings(a[i].Values)
		sort.Strings(b[i].Values)
		str1 = str1 + a[i].Name
		str2 = str2 + b[i].Name
		cnt++
		if a[i].Name != b[i].Name {
			errs++
		}
		for _, val := range a[i].Values {
			str1 = str1 + "," + val
		}
		for _, val := range b[i].Values {
			str2 = str2 + "," + val
		}
		for j, val := range a[i].Values {
			cnt++
			if j >= len(b[i].Values) || val != b[i].Values[j] {
				errs++
			}
		}
		str1 = str1 + "\n"
		str2 = str2 + "\n"
	}

	/* Write to CSV */
	err := ioutil.WriteFile("Repl"+strconv.Itoa(no)+".csv", []byte(str1), 0644)
	if err != nil {
		util.PrintErr("CHECKER", err)
	}
	err = ioutil.WriteFile("Repl"+strconv.Itoa(no+1)+".csv", []byte(str2), 0644)
	if err != nil {
		util.PrintErr("CHECKER", err)
	}
	return errs, cnt
}
