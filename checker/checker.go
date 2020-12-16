package main

import (
	"fmt"
	"os"
	"sort"

	"../util"
	"go.mongodb.org/mongo-driver/mongo"
)

// Constants
const (
	posCollection = "kvsp"
	negCollection = "kvsn"
	dCollection   = "kvsd"
)

func main() {
	var err error

	/* Parse command line arguments */
	drop := os.Args[1]
	impl := os.Args[2]
	if err != nil {
		util.PrintErr("CHECKER", "CmdLine", err)
	}

	/* Parse group member information */
	ips, ports, dbPorts, err := util.ParseGroupMembersCVS("../ports.csv", "")
	if err != nil {
		util.PrintErr("CHECKER", "GroupInfo", err)
	}
	noReplicas := len(dbPorts)

	/* Init data structures */
	dbs := make([]*mongo.Database, noReplicas)
	results := make([][]util.SRecord, noReplicas)

	/* Connect to databases */
	for i, dbPort := range dbPorts {
		dbClient, _ := util.ConnectDb("CHECKER", ips[i], dbPort)
		dbs[i] = dbClient.Database("chev")
		util.PrintMsg("CHECKER", "Connected to DB on port "+dbPort)
	}

	/* Download results and save */
	if impl == "cv" { // cv: need to download pos and neg collections
		for i, db := range dbs {
			result := util.DownloadDState(db, "CHECKER", posCollection, drop)
			util.SaveDStateToCSV(result, i, "P")
			result = util.DownloadDState(db, "CHECKER", negCollection, drop)
			util.SaveDStateToCSV(result, i, "N")
		}
	}
	if impl == "cm" { // cm: need to download dynamic collection and do lookup
		for i, db := range dbs {
			var res int
			conn := util.RPCClient("CHECKER", ips[i], ports[i])
			conn.Call("RPCExt.Lookup", util.RPCExtArgs{}, &res)
			result := util.DownloadDState(db, "CHECKER", dCollection, drop)
			util.SaveDStateToCSV(result, i, "D")
		}
	}
	for i, db := range dbs { // all: get the static collection
		results[i] = util.DownloadSState(db, "CHECKER", drop)
		util.SaveSStateToCSV(results[i], i)
	}

	/* Download the reference */
	ref := util.DownloadMainTestRef()

	/* Check Equality */
	var i1, i2 int
	var sum1, sum2 float32
	for i := 0; i < noReplicas; i++ {
		if i != noReplicas-1 {
			i1 = i     // array index of first repl
			i2 = i + 1 // array index of second repl
		} else {
			i1 = i
			i2 = 0
		}

		/* Check against each other */
		errs, cnt := testEq(results[i1], results[i2])
		avg := (cnt - errs) / cnt * 100
		msg := fmt.Sprint("Number of diffs ", i1+1, " to ", i2+1, " is ", errs,
			" and percent is ", avg)
		util.PrintMsg("CHECKER", msg)
		sum1 += avg

		/* Check against reference */
		errs, cnt = testEq(ref, results[i1])
		avg = (cnt - errs) / cnt * 100
		msg = fmt.Sprint("Number of diffs ", i1+1, " to ref is ", errs,
			" and percent is ", avg)
		util.PrintMsg("CHECKER", msg)
		sum2 += avg
	}

	/* Print overall consistencies to console */
	util.PrintMsg("CHECKER", "Overall consistency against each other is "+fmt.Sprint(sum1/float32(noReplicas)))
	util.PrintMsg("CHECKER", "Overall consistency against ref is "+fmt.Sprint(sum2/float32(noReplicas)))
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

		/* Check of keys */
		if a[i].Key != b[i].Key {
			errs++
		}
		cnt++

		/* Check values */
		for j, val := range a[i].Values {
			if j >= len(b[i].Values) || val != b[i].Values[j] {
				errs++
			}
			cnt++
		}
	}

	return errs, cnt
}
