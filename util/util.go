package util

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

// OpCode is an operation code
type OpCode int

// OpCodes
const (
	IK OpCode = iota + 1
	IV
	RK
	RV
	NO
)

// LookupOpCode translate operation code from string to op code
func LookupOpCode(opCode OpCode, noStr string) string {
	if opCode == IK {
		return "IK"
	} else if opCode == IV {
		return "IV"
	} else if opCode == RK {
		return "RK"
	} else if opCode == RV {
		return "RV"
	} else if opCode == NO {
		return "No-Op"
	} else {
		PrintErr(noStr, errors.New("LookupOpCode: unknown operation"))
		return ""
	}
}

// ParseGroupMembersCVS parses the supplied CVS group member file
func ParseGroupMembersCVS(file string, port string) ([]string, []string, []string, error) {
	// adapted from https://stackoverflow.com/questions/24999079/reading-csv-file-in-go
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	clPorts := []string{}
	dbPorts := []string{}
	ips := []string{}

	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return ips, clPorts, dbPorts, nil
		}

		/* Remove own port from results if appropriate */
		if row[1] != port {
			ips = append(ips, row[0])
			clPorts = append(clPorts, row[1])
			dbPorts = append(dbPorts, row[2])
		}
	}
}

// PrintMsg prints message to console from a replica
func PrintMsg(no string, msg string) {
	if no == "CLIENT" || no == "TESTER" {
		fmt.Println(no + ": " + msg)
	} else {
		fmt.Println("REPLICA " + no + ": " + msg)
	}
}

// PrintErr prints error to console from a replica and exits
func PrintErr(no string, err error) {
	if no == "CLIENT" || no == "TESTER" {
		fmt.Println(no+": ", err)
	} else {
		fmt.Println("REPLICA "+no+": ", err)
	}
	os.Exit(1)
}

// EmulateDelay emulates link delay in all RPC responses
func EmulateDelay(delay int) {
	if delay > 0 {
		time.Sleep(time.Duration(GetRand(delay)) * time.Millisecond)
	}
}

// GetRand generates a random number to emulate connection delays
func GetRand(no int) int {
	// https://golang.cafe/blog/golang-random-number-generator.html
	rand.Seed(time.Now().UnixNano())
	min := int(0.8 * float32(no))
	max := int(1.2 * float32(no))
	res := rand.Intn(max-min+1) + min
	return res
}

// Max returns the maximum of a and b
func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
