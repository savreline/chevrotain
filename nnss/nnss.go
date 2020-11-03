package main

import (
	"fmt"
	"net/rpc"
	"os"
	"strconv"

	"../cmrdt"
	"../util"
	"github.com/DistributedClocks/GoVector/govec"
)

var clients []*rpc.Client
var logs []*govec.GoLog

func main() {
	noReplicas, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		util.PrintErr(err)
	}
	clients = make([]*rpc.Client, noReplicas)
	logs = make([]*govec.GoLog, noReplicas)

	/* Init Replicas */
	cmrdt.InitReplica(true, "1", noReplicas, "8001", "27018")
	cmrdt.InitReplica(true, "2", noReplicas, "8002", "27019")

	/* Parse Ports */
	ports, err := util.ParseGroupMembersText("ports.txt", "")
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println("DRIVER: Connecting to ports ", ports)

	/* Make RPC Connections */
	for i, port := range ports {
		rpcChan := make(chan *rpc.Client)
		logChan := make(chan *govec.GoLog)
		go util.RPCClient(rpcChan, logChan, port, "DRIVER: ")
		clients[i] = <-rpcChan
		logs[i] = <-logChan
	}

	// fmt.Println(clients[0])
	// fmt.Println(clients[1])

	/* A few sample RPC Commands */
	var result int
	err = clients[0].Call("RPCCmd.ConnectReplica", cmrdt.ConnectArgs{Val: ""}, &result)
	if err != nil {
		util.PrintErr(err)
	}
	err = clients[1].Call("RPCCmd.ConnectReplica", cmrdt.ConnectArgs{Val: ""}, &result)
	if err != nil {
		util.PrintErr(err)
	}
	err = clients[0].Call("RPCCmd.InsertKey", cmrdt.KeyArgs{Key: "1"}, &result)
	if err != nil {
		util.PrintErr(err)
	}
	err = clients[1].Call("RPCCmd.InsertKey", cmrdt.KeyArgs{Key: "2"}, &result)
	if err != nil {
		util.PrintErr(err)
	}

	// cmrdt.TerminateReplica()
	for {
	}
}
