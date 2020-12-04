package main

import (
	"io/ioutil"
	"os"
	"strings"

	"../../util"
)

func wikiTest() {
	go loadPages("Java", 0)
	go loadPages("C--", 1)
	go loadPages("C++", 2)
}

// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
func loadPages(startPage string, no int) {
	pathHead := "../../crawler/" + startPage + "/"
	lastPage := startPage
	maxDepth := 3

	/* Connect to the Replica and Connect the Replica */
	conn := util.ConnectDriver(ports[no], t)
	var result int

	/* Init Queue */
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
		err = conn.Call("RPCExt.InsertKey", util.RPCExtArgs{Key: curPage}, &result)
		if err != nil {
			util.PrintErr("DRIVER", err)
		}

		/* Add to Queue and Insert Value */
		for j := 0; j < len(linkedPages)-1; j++ {
			queue = append(queue, linkedPages[j])

			/* Insert Value */
			err = conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: curPage, Value: linkedPages[j]}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
		}
	}

	/* Terminate */
	util.Terminate(ports[no], conn, 3)
}
