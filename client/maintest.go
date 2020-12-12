package main

import (
	"strconv"
	"sync"
	"time"

	"../util"
)

func test1(no int, noKeys int, noVals int, removes bool) {
	/* Connect to the replica and Connect the replica */
	conn := util.ConnectClient(ips[no], ports[no], timeInt)
	cnt := 0

	/* init map of latencies, associated wait group and lock */
	latencies := make(map[int]int64)
	var lock sync.Mutex
	var wg sync.WaitGroup

	/* Record starting time */
	t := time.Now().UnixNano()

	/* Inserts */
	for i := 0; i < noKeys; i++ {
		key := (no+1)*100 + i
		cnt++
		go sendCmd(strconv.Itoa(key), "", cnt, util.IK, conn, latencies, &lock, &wg)
		time.Sleep(time.Duration(delay) * time.Millisecond)

		for j := 0; j < noVals; j++ {
			val := (no+1)*1000 + j
			cnt++
			go sendCmd(strconv.Itoa(key), strconv.Itoa(val), cnt, util.IV, conn, latencies, &lock, &wg)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}

	if removes {
		/* Remove Values: remove the latter half of the values from the latter half of the keys */
		for i := noKeys / 2; i < noKeys; i++ {
			key := (no+1)*100 + i

			for j := noVals / 2; j < noVals; j++ {
				val := (no+1)*1000 + j
				cnt++
				go sendCmd(strconv.Itoa(key), strconv.Itoa(val), cnt, util.RV, conn, latencies, &lock, &wg)
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}

		/* Remove Keys: remove the last quater of the keys */
		for i := 3 * noKeys / 4; i < noKeys; i++ {
			key := (no+1)*100 + i

			cnt++
			go sendCmd(strconv.Itoa(key), "", cnt, util.RK, conn, latencies, &lock, &wg)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}

	/* Done sending commands, record time */
	delta := time.Now().UnixNano() - t
	util.PrintMsg("CLIENT", "Done Sending Calls, Waiting")
	wg.Wait()

	/* Terminate */
	util.TerminateReplica(ports[no], conn, 5)

	/* Process collected performance data */
	calcPerf(delta, cnt, no, latencies)
}
