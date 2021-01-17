package main

import (
	"testing"

	"../util"
)

// go test -run NameOfTest

/* First Insert a Number of Identical Elements with Different Ids:
* Key 100 x3
* Key 200 x2
* Key 100, Val 1000 x3
* Key 200, Val 2000 x2
* Key 300, Val 3000 (diff key)
* Key 100, Val 1001 (pre-existing key) */
var pairs = [][2]string{
	{"Keys", "100"},
	{"Keys", "100"},
	{"Keys", "100"},
	{"Keys", "200"},
	{"Keys", "200"},
	{"100", "1000"},
	{"100", "1000"},
	{"100", "1000"},
	{"200", "2000"},
	{"200", "2000"},
	{"300", "3000"},
	{"100", "1001"},
}

func TestInserts(t *testing.T) {
	db = util.ConnectLocalDb()
	/* Can we insert a key? */
	insert("Keys", "100", 1)
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "0"))
	/* Can we insert a value? */
	insert("100", "1000", 2)
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "0"))
	/* Can we insert a value without a key? */
	insert("200", "2000", 3)
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "0"))
	/* Can we insert a value into pre-existing key? */
	insert("200", "2001", 4)
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "0"))
	/* What about insering the same key again? */
	insert("Keys", "100", 5)
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "1"))
}

func TestDeletes(t *testing.T) {
	db = util.ConnectLocalDb()
	id := 1
	for _, pair := range pairs {
		insert(pair[0], pair[1], id)
		id++
	}
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "0"))

	/* Now remove those one-by-one */
	remove("Keys", "100", computeRemovalSet("Keys", "100"))
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "0"))
	remove("Keys", "200", computeRemovalSet("Keys", "200"))
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "0"))
	remove("100", "1000", computeRemovalSet("100", "1000"))
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "0"))
	remove("200", "2000", computeRemovalSet("200", "2000"))
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "0"))
	remove("300", "3000", computeRemovalSet("300", "3000"))
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "0"))
	remove("100", "1001", computeRemovalSet("100", "1001"))
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "1"))
}

func TestLookup(t *testing.T) {
	db = util.ConnectLocalDb()
	id := 1
	for _, pair := range pairs {
		insert(pair[0], pair[1], id)
		id++
	}
	var res int
	rpcext := new(RPCExt)
	rpcext.Lookup(&util.RPCExtArgs{}, &res)
	util.PrintSState(util.DownloadSState(db, "TESTER", "1"))
	util.PrintDState(util.DownloadDState(db, "TESTER", dCollection, "1"))
}
