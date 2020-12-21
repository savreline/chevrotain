package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"../util"
)

// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
func wikiTest(startPage string, no int) {
	defer wgMain.Done()
	pathHead := "../crawler/" + startPage + "/"
	lastPage := startPage
	maxDepth := 3

	/* Connect to the Replica and Connect the Replica */
	conn := util.ConnectClient(ips[no], ports[no], timeInt)

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
		go sendCmd(curPage, "", util.IK, conn)
		time.Sleep(time.Duration(delay) * time.Microsecond)

		/* Add to Queue and Insert Value */
		for j := 0; j < len(linkedPages)-1; j++ {
			queue = append(queue, linkedPages[j])

			/* Insert Value */
			go sendCmd(curPage, linkedPages[j], util.IV, conn)
			time.Sleep(time.Duration(delay) * time.Microsecond)
		}
	}

	/* Done sending commands, record time */
	delta := time.Now().UnixNano() - t
	util.PrintMsg("CLIENT", "Done Sending Calls, Waiting, Delta: "+fmt.Sprint(delta/1000000))
	wg.Wait()

	/* Terminate */
	util.TerminateReplica(ports[no], conn, 3)

	/* Process collected performance data */
	calcPerf()
}
