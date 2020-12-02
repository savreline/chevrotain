package main

import (
	"fmt"
	"net/rpc"
	"strconv"
	"time"

	"../../util"
	"github.com/savreline/GoVector/govec"
	"github.com/savreline/GoVector/govec/vclock"
)

// StateArgs are the arguments to the MergeState RPC
type StateArgs struct {
	PosState, NegState []util.CvRecord
	No                 string
	Timestamp          vclock.VClock
}

// TickArgs are the arguments passed with RPC calls that update the tick
type TickArgs struct {
	Tick      int
	No        string
	Timestamp vclock.VClock
}

// MergeState merges incoming state with the current state at the replica
func (t *RPCInt) MergeState(args *StateArgs, reply *int) error {
	msg := "Recv State from " + args.No
	logger.MergeIncomingClock(msg, args.Timestamp, govec.GetDefaultLogOptions().Priority)
	mergeCollection(args.PosState, posCollection)
	mergeCollection(args.NegState, negCollection)
	return nil
}

// TickClock accepts incoming tick of the other replica and sends out its own curTick
func (t *RPCInt) TickClock(args *TickArgs, reply *int) error {
	msg := "Recv Tick from " + args.No
	logger.MergeIncomingClock(msg, args.Timestamp, govec.GetDefaultLogOptions().Priority)
	*reply = curTick
	return nil
}

// GetCurrentSnapShot merges the positive and negative datasets into one collection
func GetCurrentSnapShot() error {
	fmt.Println("Garbage Collecting")
	mergeCollections()
	return nil
}

// broadcasts state to all other replicas
func sendState() {
	<-chanState
	for {
		time.Sleep(time.Duration(timeInt) * time.Millisecond)
		posState := util.DownloadCvState(db.Collection(posCollection), "0")
		negState := util.DownloadCvState(db.Collection(negCollection), "0")
		logger.StartBroadcast("OUT "+noStr, govec.GetDefaultLogOptions())
		broadcast(StateArgs{PosState: posState,
			NegState:  negState,
			No:        noStr,
			Timestamp: logger.GetCurrentVC().Copy()},
			TickArgs{},
			"RPCInt.MergeState",
			"RPC Sending State")
		logger.StopBroadcast()
	}
}

// broadcasts tick to all other replicas
func sendTick() {
	<-chanTick
	for {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		var str string
		curTick, str = util.UpdateCurTick(ticks, noReplicas, curTick)
		eLog = eLog + str
		logger.StartBroadcast("TICK OUT "+noStr, govec.GetDefaultLogOptions())
		calls, results := broadcast(StateArgs{},
			TickArgs{Tick: curTick,
				No:        noStr,
				Timestamp: logger.GetCurrentVC().Copy()},
			"RPCInt.TickClock",
			"RPC Sending Tick")
		logger.StopBroadcast()
		waitForBroadcastToFinish(calls)
		updateLocalTick(results)
	}
}

func runGC() {
	<-chanGC
	for {
		time.Sleep(time.Duration(10000) * time.Millisecond)
		GetCurrentSnapShot()
	}
}

// broadcasts state to all other replicas
func broadcast(state StateArgs, tick TickArgs, name string, msg string) ([]*rpc.Call, []int) {
	var results = make([]int, len(conns))
	var destNo int
	var err error
	var flag = false
	var calls = make([]*rpc.Call, len(conns))

	for i, client := range conns {
		if i == no {
			flag = true
		}

		if client != nil {
			if flag {
				destNo = i + 2
			} else {
				destNo = i + 1
			}
			fmt.Println(msg, no, "->", destNo)
			if name == "RPCInt.MergeState" {
				calls[i] = client.Go(name, state, &results[i], nil)
			} else {
				calls[i] = client.Go(name, tick, &results[i], nil)
			}
			if err != nil {
				util.PrintErr(noStr, err)
			}
		}
	}
	return calls, results
}

func addTicks(timestamp vclock.VClock) {
	for i := 0; i < noReplicas; i++ {
		ticks[i] = append(ticks[i], int(timestamp["R"+strconv.Itoa(i+1)]))
	}
}

func updateLocalTick(results []int) {
	curLocalTick = curTick
	for i := 0; i < len(results)-1; i++ {
		if results[i] < curLocalTick {
			curLocalTick = results[i]
		}
	}
}

func waitForBroadcastToFinish(calls []*rpc.Call) {
	for _, call := range calls {
		if call != nil {
			<-call.Done
		}
	}
}
