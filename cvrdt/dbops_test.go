package main

import (
	"fmt"
	"testing"

	"../util"
)

const (
	posCollection1 = "kvsp1"
	negCollection1 = "kvsn1"
	posCollection2 = "kvsp2"
	negCollection2 = "kvsn2"
)

var inserts1 = [][2]string{
	{"Keys", "100"},
	{"Keys", "100"},
	{"Keys", "200"},
	{"100", "1000"},
	{"100", "1000"},
	{"200", "2000"},
	{"300", "3000"},
	{"100", "1001"},
}

var removes1 = [][2]string{
	{"Keys", "100"},
	{"100", "1000"},
	{"200", "2000"},
	{"100", "1001"},
}

var inserts2 = [][2]string{
	{"Keys", "100"},
	{"100", "1001"},
}

func TestStateMerges(t *testing.T) {
	db = util.ConnectLocalDb()
	noStr = "1"
	clock = 1
	for i, pair := range inserts1 {
		if i%2 == 0 {
			insertLocalRecord(pair[0], pair[1], posCollection1, nil)
		} else {
			insertLocalRecord(pair[0], pair[1], posCollection2, nil)
		}
	}
	for i, pair := range removes1 {
		if i%2 == 0 {
			insertLocalRecord(pair[0], pair[1], negCollection1, nil)
		} else {
			insertLocalRecord(pair[0], pair[1], negCollection2, nil)
		}
	}
	for i, pair := range inserts2 {
		if i%2 == 0 {
			insertLocalRecord(pair[0], pair[1], posCollection1, nil)
		} else {
			insertLocalRecord(pair[0], pair[1], posCollection2, nil)
		}
	}
	posState1 := util.DownloadDState(db.Collection(posCollection1), "TESTER", "0")
	negState1 := util.DownloadDState(db.Collection(negCollection1), "TESTER", "0")
	posState2 := util.DownloadDState(db.Collection(posCollection2), "TESTER", "1")
	negState2 := util.DownloadDState(db.Collection(negCollection2), "TESTER", "1")
	mergeState(posState2, posCollection1)
	mergeState(negState2, negCollection1)
	util.PrintDState(posState1)
	util.PrintDState(posState2)
	util.PrintDState(util.DownloadDState(db.Collection(posCollection1), "TESTER", "1"))
	util.PrintDState(negState1)
	util.PrintDState(negState2)
	util.PrintDState(util.DownloadDState(db.Collection(negCollection1), "TESTER", "1"))
}

func TestColMerges1A(t *testing.T) {
	db = util.ConnectLocalDb()
	noStr = "1"
	clock = 1
	curSafeTick = 7

	/* All records distinct, all just in the positive collection */
	insertLocalRecord("Keys", "100", posCollection, nil)
	insertLocalRecord("Keys", "200", posCollection, nil)
	insertLocalRecord("Keys", "300", posCollection, nil)
	insertLocalRecord("100", "1000", posCollection, nil)
	insertLocalRecord("100", "1001", posCollection, nil)
	insertLocalRecord("200", "2000", posCollection, nil)
	mergeCollections()
	util.PrintDState(util.DownloadDState(db.Collection(posCollection), "TESTER", "1"))
	util.PrintSState(util.DownloadSState(db.Collection(sCollection), "TESTER", "1"))
}

func TestColMerges1B(t *testing.T) {
	db = util.ConnectLocalDb()
	noStr = "1"
	clock = 1
	curSafeTick = 8

	/* All records distinct, one of the KEYS will appear in the negative collection */
	insertLocalRecord("Keys", "100", posCollection, nil)
	insertLocalRecord("Keys", "200", posCollection, nil)
	insertLocalRecord("Keys", "300", posCollection, nil)
	insertLocalRecord("100", "1000", posCollection, nil)
	insertLocalRecord("100", "1001", posCollection, nil)
	insertLocalRecord("200", "2000", posCollection, nil)
	insertLocalRecord("Keys", "300", negCollection, nil)
	mergeCollections()
	util.PrintDState(util.DownloadDState(db.Collection(posCollection), "TESTER", "1"))
	util.PrintDState(util.DownloadDState(db.Collection(negCollection), "TESTER", "1"))
	util.PrintSState(util.DownloadSState(db.Collection(sCollection), "TESTER", "1"))
}

func TestColMerges1C(t *testing.T) {
	db = util.ConnectLocalDb()
	noStr = "1"
	clock = 1
	curSafeTick = 9

	/* All records distinct, one of the KEYS and VALUES will appear in the negative collection */
	insertLocalRecord("Keys", "100", posCollection, nil)
	insertLocalRecord("Keys", "200", posCollection, nil)
	insertLocalRecord("Keys", "300", posCollection, nil)
	insertLocalRecord("100", "1000", posCollection, nil)
	insertLocalRecord("100", "1001", posCollection, nil)
	insertLocalRecord("200", "2000", posCollection, nil)
	insertLocalRecord("Keys", "300", negCollection, nil)
	insertLocalRecord("100", "1001", negCollection, nil)
	mergeCollections()
	util.PrintDState(util.DownloadDState(db.Collection(posCollection), "TESTER", "1"))
	util.PrintDState(util.DownloadDState(db.Collection(negCollection), "TESTER", "1"))
	util.PrintSState(util.DownloadSState(db.Collection(sCollection), "TESTER", "1"))
}

func TestColMerges1D(t *testing.T) {
	db = util.ConnectLocalDb()
	noStr = "1"
	clock = 1
	curSafeTick = 20

	/* Extension of the previous scenarios */
	insertLocalRecord("Keys", "100", posCollection, nil)
	insertLocalRecord("Keys", "200", posCollection, nil)
	insertLocalRecord("Keys", "300", posCollection, nil)
	insertLocalRecord("Keys", "400", posCollection, nil)
	insertLocalRecord("100", "1000", posCollection, nil)
	insertLocalRecord("100", "1001", posCollection, nil)
	insertLocalRecord("100", "1002", posCollection, nil)
	insertLocalRecord("200", "2000", posCollection, nil)
	insertLocalRecord("200", "2001", posCollection, nil)
	insertLocalRecord("300", "3000", posCollection, nil)
	insertLocalRecord("Keys", "400", negCollection, nil)
	insertLocalRecord("100", "1002", negCollection, nil)
	insertLocalRecord("100", "1001", negCollection, nil)
	insertLocalRecord("200", "2001", negCollection, nil)
	insertLocalRecord("300", "3000", negCollection, nil)
	insertLocalRecord("Keys", "300", negCollection, nil)
	mergeCollections()
	util.PrintDState(util.DownloadDState(db.Collection(posCollection), "TESTER", "1"))
	util.PrintDState(util.DownloadDState(db.Collection(negCollection), "TESTER", "1"))
	util.PrintSState(util.DownloadSState(db.Collection(sCollection), "TESTER", "0"))
}

func TestColMerges2(t *testing.T) {
	db = util.ConnectLocalDb()
	noStr = "1"
	clock = 1
	curSafeTick = 5

	/* Key 100 will appear later in the negative collection
	* Key 200 will appear later in the positive collection */
	insertLocalRecord("Keys", "100", posCollection, nil)
	insertLocalRecord("Keys", "100", negCollection, nil)
	insertLocalRecord("Keys", "200", negCollection, nil)
	insertLocalRecord("Keys", "200", posCollection, nil)
	util.PrintDState(util.DownloadDState(db.Collection(posCollection), "TESTER", "0"))
	util.PrintDState(util.DownloadDState(db.Collection(negCollection), "TESTER", "0"))
	mergeCollections()
	util.DownloadDState(db.Collection(posCollection), "TESTER", "1")
	util.DownloadDState(db.Collection(negCollection), "TESTER", "1")
	util.PrintSState(util.DownloadSState(db.Collection(sCollection), "TESTER", "0"))
}

func TestColMerges3(t *testing.T) {
	db = util.ConnectLocalDb()
	noStr = "1"
	clock = 1
	curSafeTick = 5

	/* Key 100 will appear latest in the positive collection */
	insertLocalRecord("Keys", "100", posCollection, nil)
	insertLocalRecord("Keys", "100", negCollection, nil)
	insertLocalRecord("Keys", "100", negCollection, nil)
	insertLocalRecord("Keys", "100", posCollection, nil)
	mergeCollections()
	util.PrintSState(util.DownloadSState(db.Collection(sCollection), "TESTER", "1"))
}

func TestColMerges4(t *testing.T) {
	db = util.ConnectLocalDb()
	noStr = "1"
	clock = 1
	curSafeTick = 10

	/* Value 1000 will appear later in the positive collection
	* Value 2000 will appear later in the negative collection */
	insertLocalRecord("Keys", "100", posCollection, nil)
	insertLocalRecord("100", "1000", posCollection, nil)
	insertLocalRecord("100", "1000", negCollection, nil)
	insertLocalRecord("Keys", "200", posCollection, nil)
	insertLocalRecord("200", "2000", negCollection, nil)
	insertLocalRecord("200", "2000", posCollection, nil)
	util.PrintDState(util.DownloadDState(db.Collection(posCollection), "TESTER", "0"))
	util.PrintDState(util.DownloadDState(db.Collection(negCollection), "TESTER", "0"))
	mergeCollections()
	util.DownloadDState(db.Collection(posCollection), "TESTER", "1")
	util.DownloadDState(db.Collection(negCollection), "TESTER", "1")
	util.PrintSState(util.DownloadSState(db.Collection(sCollection), "TESTER", "1"))
}

func TestColMerges5(t *testing.T) {
	db = util.ConnectLocalDb()
	noStr = "1"
	clock = 1
	curSafeTick = 10

	/* Value 1000 will appear latest in the positive collection */
	insertLocalRecord("Keys", "100", posCollection, nil)
	insertLocalRecord("100", "1000", posCollection, nil)
	insertLocalRecord("100", "1000", negCollection, nil)
	insertLocalRecord("100", "1000", negCollection, nil)
	insertLocalRecord("100", "1000", posCollection, nil)
	util.PrintDState(util.DownloadDState(db.Collection(posCollection), "TESTER", "0"))
	util.PrintDState(util.DownloadDState(db.Collection(negCollection), "TESTER", "0"))
	mergeCollections()
	util.DownloadDState(db.Collection(posCollection), "TESTER", "1")
	util.DownloadDState(db.Collection(negCollection), "TESTER", "1")
	util.PrintSState(util.DownloadSState(db.Collection(sCollection), "TESTER", "1"))
}

func TestArrayShifts(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}

	j := 0
	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
		arr[j] = arr[i]
		j++
	}

	arr = arr[:j]
	fmt.Println(arr)
}
