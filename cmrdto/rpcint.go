package main

import (
	"fmt"

	"golang.org/x/sys/windows"

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

	if !ready {
		/* Make a channel to communicate on with this RPC call */
		channel := make(chan bool, 100000)

		/* Add the channel to the pool */
		lock.Lock()
		chans = append(chans, channel)
		lock.Unlock()

		/* Wait for the correct clock */
		for i := 0; ; i++ {
			if verbose > 0 {
				eLog = eLog + fmt.Sprint(windows.GetCurrentThreadId()) +
					fmt.Sprint(" IC ", args.Clock) +
					fmt.Sprint(" I: ", i) + "\n"
			}
			<-channel
			myClock = logger.GetCurrentVCSafe()
			ready = myClock.CompareBroadcastClock(args.Clock)
			if ready {
				break
			}
		}

		if verbose > 0 {
			iLog = iLog + fmt.Sprint(windows.GetCurrentThreadId()) +
				fmt.Sprint(" IC ", args.Clock) +
				" Out of loop, waiting for lock\n"
		}

		/* Remove the channel from the pool */
		lock.Lock()
		close(channel)
		removeChan(channel)
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

// remove a channel from the pool
// https://stackoverflow.com/questions/28699485/remove-elements-in-slice
func removeChan(channel chan bool) {
	for i := 0; i < len(chans); i++ {
		if channel == chans[i] {
			chans = append(chans[:i], chans[i+1:]...)
			i--
			break
		}
	}
}
