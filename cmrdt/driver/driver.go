package main

import (
	"io/ioutil"
	"strconv"
	"strings"

	"../../util"
	"github.com/DistributedClocks/GoVector/govec"
)

var ports []string
var logger *govec.GoLog

func main() {
	/* Parse Group Membership */
	var err error
	ports, _, err = util.ParseGroupMembersCVS("ports.csv", "")
	if err != nil {
		util.PrintErr(err)
	}
	noReplicas := len(ports)

	/* Init Cloks */
	logger = govec.InitGoVector("Drv", "Drv", govec.GetDefaultConfig())

	/* Tests */
	for i := 0; i < noReplicas; i++ {
		// go simpleTest(i)
	}
	wikiTest()

	for {
	}
}

// simpleTest
func simpleTest(no int) {
	/* Connect to the Replica and Connect the Replica */
	var result int
	conn := util.RPCClient(logger, ports[no], "DRIVER: ")
	err := conn.Call("RPCExt.ConnectReplica", util.ConnectArgs{}, &result)
	if err != nil {
		util.PrintErr(err)
	}

	/* Inserts */
	for i := 0; i < 10; i++ {
		num := (i + 1) * (no + 1)
		conn.Call("RPCExt.InsertKey", util.KeyArgs{Key: strconv.Itoa(num)}, &result)
	}
}

// wikiTest
func wikiTest() {
	go loadPages("Java", 0)
	go loadPages("C--", 1)
	go loadPages("C++", 2)
}

func loadPages(startPage string, no int) {
	/* Connect to the Replica and Connect the Replica */
	var result int
	conn := util.RPCClient(logger, ports[no], "DRIVER: ")
	err := conn.Call("RPCExt.ConnectReplica", util.ConnectArgs{}, &result)
	if err != nil {
		util.PrintErr(err)
	}

	/* Init Queue */
	var queue []string
	queue = append(queue, startPage)

	/* BFS */
	for len(queue) > 0 {
		/* Pop off queue */
		curPage := queue[0]
		queue = queue[1:]

		/* Insert Key */
		err := conn.Call("RPCExt.InsertKey", util.KeyArgs{Key: curPage}, &result)
		if err != nil {
			util.PrintErr(err)
		}

		/* Read the file */
		curPage = strings.Replace(curPage, "_", " ", -1)
		body, err := ioutil.ReadFile("../../crawler/" + startPage + "/" + curPage + ".link")
		if err != nil {
			break
		}

		/* Add the linked files to queue */
		linkedPages := strings.Split(string(body[:]), "\n")
		for _, page := range linkedPages {
			if page != "" {
				queue = append(queue, page)
			}
		}
	}
}
