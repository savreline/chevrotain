package main

import (
	"fmt"
	"net/rpc"
	"strconv"
	"time"

	"../util"
)

func maintest(no int) {
	defer wgMain.Done()

	/* Connect to the replica and Connect the replica */
	var conn *rpc.Client
	if !mongotest {
		conn = util.ConnectClient(ips[no], ports[no], timeInt)
	}

	/* Record starting time */
	t := time.Now().UnixNano()

	/* Inserts */
	for i := no * noPerRepl; i < (no+1)*noPerRepl; i++ {
		go sendCmd(strconv.Itoa(i+100), "", util.IK, conn)
		time.Sleep(time.Duration(delay) * time.Microsecond)

		for j := 0; j < noVals; j++ {
			go sendCmd(strconv.Itoa(i+100), strconv.Itoa(j+1000), util.IV, conn)
			time.Sleep(time.Duration(delay) * time.Microsecond)
		}
	}

	if removes {
		/* Remove Values */
		for i := no * noPerRepl; i < (no+1)*noPerRepl; i++ {
			if i%2 != 0 {
				continue
			}

			for j := 0; j < noVals/2; j++ {
				go sendCmd(strconv.Itoa(i+100), strconv.Itoa(j+1000), util.RV, conn)
				time.Sleep(time.Duration(delay) * time.Microsecond)
			}
		}

		/* Remove Keys */
		for i := no * noPerRepl; i < (no+1)*noPerRepl; i++ {
			if i%4 != 0 {
				continue
			}

			go sendCmd(strconv.Itoa(i+100), "", util.RK, conn)
			time.Sleep(time.Duration(delay) * time.Microsecond)
		}
	}

	/* Done sending commands, record time */
	delta := time.Now().UnixNano() - t
	util.PrintMsg("CLIENT", "Done Sending Calls, Waiting, Delta: "+fmt.Sprint(delta/1000000))
	wg.Wait()

	/* Terminate */
	if !mongotest && term {
		util.TerminateReplica(ports[no], conn, 5)
	}
}
