package main

import (
	"fmt"

	"../util"
	"github.com/savreline/GoVector/govec"
)

// ProcessIntCall processes the broadcasts coming from other replicas
func (t *RPCInt) ProcessIntCall(args *BroadcastArgs, reply *int) error {
	/* Hold the call until its turn is up */
	waitForTurn(args)

	/* Execute effect-update */
	effectUpdates(args.Key, args.Value, args.Ids, args.OpType)

	/* Finish */
	util.EmulateDelay(delay)
	return nil
}

// wait for the correct turn to process the call
func waitForTurn(args *BroadcastArgs) {
	/* Check if this RPC call needs to wait */
	myClock := logger.GetCurrentVCSafe()
	ready := myClock.CompareBroadcastClock(args.Clock)
	if verbose {
		eLog = eLog + fmt.Sprint("K:", args.Key) + fmt.Sprint(" V:", args.Value) +
			fmt.Sprint(" My clock ", myClock) + fmt.Sprint(" Incoming clock ", args.Clock) +
			fmt.Sprint(" Comparison: ", ready) + "\n"
	}

	if ready == false {
		/* Make a channel to communicate on with this RPC call */
		channel := make(chan bool, 100)

		/* Add the channel to the pool */
		lock.Lock()
		chans[channel] = channel
		lock.Unlock()

		/* Wait for the correct clock */
		for i := 0; ; i++ {
			<-channel
			myClock = logger.GetCurrentVCSafe()
			ready = myClock.CompareBroadcastClock(args.Clock)
			if verbose {
				iLog = iLog + fmt.Sprint("K:", args.Key) + fmt.Sprint(" V:", args.Value) +
					fmt.Sprint(" My clock: ", myClock) + fmt.Sprint(" Incoming clock: ", args.Clock) +
					fmt.Sprint(" Comparison: ", ready) + fmt.Sprint(" Iteration: ", i) + "\n"
			}
			if ready {
				break
			}
		}

		/* Remove the channel from the pool */
		lock.Lock()
		delete(chans, channel)
		lock.Unlock()
	}

	/* Merge the clock */
	msg := "IN " + util.LookupOpCode(args.OpType, noStr) + " on R " + args.SrcPid +
		" K:" + args.Key + " V:" + args.Value
	logger.MergeIncomingClock(msg, args.Clock, govec.GetDefaultLogOptions().Priority)
	broadcastNewMerge()
}

// alert all waiting calls that a new clock has been merged in
func broadcastNewMerge() {
	lock.Lock()
	for _, channel := range chans {
		channel <- true
	}
	lock.Unlock()
}
