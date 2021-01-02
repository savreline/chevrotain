package main

import (
	"fmt"
	"strconv"
	"time"

	"../util"
)

func quicktest(no int) {
	defer wgMain.Done()

	/* Connect to the replica and Connect the replica */
	conn := util.ConnectClient(ips[no], ports[no], timeInt)

	/* Record starting time */
	t := time.Now().UnixNano()

	/* Inserts */
	for i := 0; i < noKeys; i++ {
		go sendCmd(strconv.Itoa((no+1)*100+i), "", util.IK, conn)
		time.Sleep(time.Duration(delay) * time.Microsecond)

		for j := 0; j < noVals; j++ {
			go sendCmd(strconv.Itoa((no+1)*100+i), strconv.Itoa((no+1)*1000+j), util.IV, conn)
			time.Sleep(time.Duration(delay) * time.Microsecond)
		}
	}

	if removes {
		/* Remove Values: remove the latter half of the values from the latter half of the keys */
		for i := noKeys / 2; i < noKeys; i++ {
			for j := noVals / 2; j < noVals; j++ {
				go sendCmd(strconv.Itoa((no+1)*100+i), strconv.Itoa((no+1)*1000+j), util.RV, conn)
				time.Sleep(time.Duration(delay) * time.Microsecond)
			}
		}

		/* Remove Keys: remove the last quater of the keys */
		for i := 3 * noKeys / 4; i < noKeys; i++ {
			go sendCmd(strconv.Itoa((no+1)*100+i), "", util.RK, conn)
			time.Sleep(time.Duration(delay) * time.Microsecond)
		}
	}

	/* Done sending commands, record time */
	delta := time.Now().UnixNano() - t
	util.PrintMsg("CLIENT", "Done Sending Calls, Waiting, Delta: "+fmt.Sprint(delta/1000000))
	wg.Wait()

	/* Terminate */
	util.TerminateReplica(ports[no], conn, 5)
}
