package main

import (
	"net/rpc"

	"../../util"
	"../cvrdt"
	"github.com/DistributedClocks/GoVector/govec"
)

var clients []*rpc.Client

func main() {
	/* Parse Group Membership */
	clPorts, dbPorts, err := util.ParseGroupMembersCVS("ports.csv", "")
	if err != nil {
		util.PrintErr(err)
	}
	noReplicas := len(clPorts)
	clients = make([]*rpc.Client, noReplicas)

	/* Init Cloks */
	logger := govec.InitGoVector("Drv", "Drv", govec.GetDefaultConfig())

	/* Init Replicas */
	cvrdt.Init(noReplicas, true, 600)
	for i := 0; i < noReplicas; i++ {
		cvrdt.InitReplica(i, clPorts[i], dbPorts[i])
	}

	/* Make RPC Connections */
	for i, port := range clPorts {
		channel := make(chan *rpc.Client)
		go util.RPCClient(channel, logger, port, "DRIVER: ")
		clients[i] = <-channel
	}

	/* A few sample RPC Commands */
	var result int
	go func() {
		err = clients[0].Call("RPCCmd.ConnectReplica", cvrdt.ConnectArgs{No: 0}, &result)
		if err != nil {
			util.PrintErr(err)
		}
		err = clients[0].Call("RPCCmd.InsertKey", cvrdt.KeyArgs{No: 0, Key: "1"}, &result)
		if err != nil {
			util.PrintErr(err)
		}
	}()
	go func() {
		err = clients[1].Call("RPCCmd.ConnectReplica", cvrdt.ConnectArgs{No: 1}, &result)
		if err != nil {
			util.PrintErr(err)
		}
		err = clients[1].Call("RPCCmd.InsertKey", cvrdt.KeyArgs{No: 1, Key: "2"}, &result)
		if err != nil {
			util.PrintErr(err)
		}
	}()
	go func() {
		err = clients[2].Call("RPCCmd.ConnectReplica", cvrdt.ConnectArgs{No: 2}, &result)
		if err != nil {
			util.PrintErr(err)
		}
		err = clients[2].Call("RPCCmd.InsertKey", cvrdt.KeyArgs{No: 2, Key: "3"}, &result)
		if err != nil {
			util.PrintErr(err)
		}
	}()

	// cmrdt.TerminateReplica()
	for {
	}
}
