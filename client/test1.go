package main

import (
	"strconv"
	"sync"
	"time"

	"../util"
)

func test1(no int, noKeys int, noVals int) {
	/* Connect to the replica and Connect the replica */
	conn := util.ConnectClient(ports[no], timeInt)
	cnt := 0

	/* init map of latencies and associated wait group */
	latencies := make(map[int]int64)
	var wg sync.WaitGroup

	/* Record starting time */
	t := time.Now().UnixNano()

	/* Inserts */
	for i := 0; i < noKeys; i++ {
		key := (no+1)*100 + i
		cnt++
		go sendCmd(strconv.Itoa(key), "", cnt, util.IK, conn, latencies, &wg)
		time.Sleep(time.Duration(delay) * time.Millisecond)

		for j := 0; j < noVals; j++ {
			val := (no+1)*1000 + j
			cnt++
			go sendCmd(strconv.Itoa(key), strconv.Itoa(val), cnt, util.IV, conn, latencies, &wg)
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
