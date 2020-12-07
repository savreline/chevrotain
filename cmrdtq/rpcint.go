package main

import (
	"../util"
	"github.com/savreline/GoVector/govec"
)

// ProcessIntCall processes an internal RPC call (i.e. an incoming broadcast)
func (t *RPCInt) ProcessIntCall(args *OpNode, reply *int) error {
	/* Add the operation to local queue */
	addToQueue(*args)

	/* Merge the incoming clock */
	var msg string
	if args.Value == "" {
		msg = "IN InsKey " + args.Key + " from " + args.SrcPid
	} else {
		msg = "IN InsVal " + args.Key + ":" + args.Value + " from " + args.SrcPid
	}
	logger.MergeIncomingClock(msg, args.Timestamp, govec.GetDefaultLogOptions().Priority)

	/* Finish */
	util.EmulateDelay(delay)
	return nil
}
