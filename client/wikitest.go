package main

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"../util"
)

// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
func wikiTest(startPage string, no int) {
	pathHead := "../crawler/" + startPage + "/"
	lastPage := startPage
	maxDepth := 3

	/* Connect to the Replica and Connect the Replica */
	conn := util.ConnectClient(ips[no], ports[no], timeInt)
	cnt := 0

	/* init map of latencies, associated wait group and lock */
	latencies := make(map[int]int64)
	var lock sync.Mutex
	var wg sync.WaitGroup

	/* Record starting time */
	t := time.Now().UnixNano()

	/* Init BSF queue */
	var queue []string
	queue = append(queue, startPage)

	/* BFS */
	i := 0
	for len(queue) > 0 && i < maxDepth {
		if queue[0] == lastPage {
			i++
		}

		/* Pop off queue */
		curPage := queue[0]
		queue = queue[1:]

		/* Read the file */
		curPageSp := strings.Replace(curPage, "_", " ", -1)
		body, err := ioutil.ReadFile(pathHead + curPageSp + ".link")
		if err != nil {
			continue
		}

		/* Add the linked files to queue: check to make sure last file actually exists */
		linkedPages := strings.Split(string(body[:]), "\n")
		if curPage == lastPage && i < maxDepth {
			m := 2
			for {
				lastPage = linkedPages[len(linkedPages)-m]
				lastPageSp := strings.Replace(lastPage, "_", " ", -1)
				if _, err := os.Stat(pathHead + lastPageSp + ".link"); err == nil {
					break
				}
				m++
			}
		}

		/* Insert Key */
		cnt++
		go sendCmd(curPage, "", cnt, util.IK, conn, latencies, &lock, &wg)
		time.Sleep(time.Duration(delay) * time.Millisecond)

		/* Add to Queue and Insert Value */
		for j := 0; j < len(linkedPages)-1; j++ {
			queue = append(queue, linkedPages[j])

			/* Insert Value */
			cnt++
			go sendCmd(curPage, linkedPages[j], cnt, util.IV, conn, latencies, &lock, &wg)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}

	/* Done sending commands, record time */
	delta := time.Now().UnixNano() - t
	util.PrintMsg("CLIENT", "Done Sending Calls, Waiting")
	wg.Wait()

	/* Terminate */
	util.TerminateReplica(ports[no], conn, 3)

	/* Process collected performance data */
	calcPerf(delta, cnt, no, latencies)
}
