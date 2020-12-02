package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"

	"../util"
	"go.mongodb.org/mongo-driver/mongo"
)

var delay int
var verbose = true
var db *mongo.Database
var latencies map[int]int64 = make(map[int]int64)

func main() {
	var err error
	noStr := os.Args[1]
	dbPort := os.Args[2]
	delay, err = strconv.Atoi(os.Args[3])
	if err != nil {
		util.PrintErr("DRIVER", err)
	}

	dbClient, _ := util.ConnectDb(noStr, dbPort)
	db = dbClient.Database("chev")
	fmt.Println("Connected to DB on " + dbPort)

	simpleTest(50, 20)
}

func simpleTest(noKeys int, noVals int) {
	k := 1
	insertKey("0")

	for i := 0; i < noKeys; i++ {
		key := strconv.Itoa(100 + i)
		t := time.Now().UnixNano()
		insertKey(key)
		latencies[k] = time.Now().UnixNano() - t
		time.Sleep(time.Duration(delay) * time.Millisecond)
		k++
		// }

		// for i := 0; i < noKeys; i++ {
		// 	key := (no+1)*100 + i
		for j := 0; j < noVals; j++ {
			val := strconv.Itoa(1000 + k)
			t := time.Now().UnixNano()
			insertValue(key, val)
			latencies[k] = time.Now().UnixNano() - t
			time.Sleep(time.Duration(delay) * time.Millisecond)
			k++
		}
	}

	/* Write Latencies to CSV */
	var str string
	var sum int64
	var keys []int
	for key := range latencies {
		keys = append(keys, key)
		sum += latencies[key]
	}
	sort.Ints(keys)
	for _, key := range keys {
		str = str + strconv.Itoa(key) + "," + strconv.FormatInt(latencies[key], 10) + "\n"
	}
	err := ioutil.WriteFile("Latencies.csv", []byte(str), 0644)
	if err != nil {
		util.PrintErr("DRIVER", err)
	}
	avg := float32(sum) / 1000000 / float32(k)
	util.PrintMsg("DRIVER:", "Average Latency is (ms):"+fmt.Sprint(avg))
}
