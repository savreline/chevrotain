package main

import (
	"io/ioutil"
	"net/rpc"
	"strconv"
	"strings"

	"../../util"
	"../cmrdt"
	"github.com/DistributedClocks/GoVector/govec"
)

var clients []*rpc.Client
var noReplicas int

func main() {
	/* Parse Group Membership */
	// clPorts, _, err := util.ParseGroupMembersCVS("ports.csv", "")
	clPorts, dbPorts, err := util.ParseGroupMembersCVS("ports.csv", "")
	if err != nil {
		util.PrintErr(err)
	}
	noReplicas = len(clPorts)
	clients = make([]*rpc.Client, noReplicas)

	/* Init Cloks */
	logger := govec.InitGoVector("Drv", "Drv", govec.GetDefaultConfig())

	/* Init Replicas */
	cmrdt.Init(noReplicas)
	for i := 0; i < noReplicas; i++ {
		cmrdt.InitReplica(true, i, clPorts[i], dbPorts[i])
	}

	/* Make RPC Connections */
	for i, port := range clPorts {
		clients[i] = util.RPCClient(logger, port, "DRIVER: ")
	}

	// simpleTest()
	wikiTest()

	// cmrdt.TerminateReplica()
	for {
	}
}

// wikiTest
func wikiTest() {
	loadPages("Java", 0)
	loadPages("C--", 1)
	loadPages("C++", 2)
}

func loadPages(startPage string, no int) {
	/* Connect to Replica */
	var result int
	err := clients[no].Call("RPCExt.ConnectReplica", cmrdt.ConnectArgs{No: no}, &result)
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
		err := clients[no].Call("RPCExt.InsertKey", cmrdt.KeyArgs{No: no, Key: curPage}, &result)
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

/* A few sample RPC Commands */
func simpleTest() {
	for i := 0; i < noReplicas; i++ {
		go initInsert(i)
	}
}

func initInsert(no int) {
	str := strconv.Itoa(no + 1)
	var result int
	err := clients[no].Call("RPCExt.ConnectReplica", cmrdt.ConnectArgs{No: no}, &result)
	if err != nil {
		util.PrintErr(err)
	}
	err = clients[no].Call("RPCExt.InsertKey", cmrdt.KeyArgs{No: no, Key: str}, &result)
	if err != nil {
		util.PrintErr(err)
	}
}
