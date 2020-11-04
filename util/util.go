package util

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DistributedClocks/GoVector/govec"
	"github.com/DistributedClocks/GoVector/govec/vrpc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect to MongoDB on the given port, as per https://www.mongodb.com/golang
func Connect(port string) (*mongo.Client, context.Context) {
	urlString := "mongodb://localhost:" + port + "/"

	client, err := mongo.NewClient(options.Client().ApplyURI(urlString))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client, ctx
}

// ParseGroupMembersCVS parses the supplied CVS group member file
func ParseGroupMembersCVS(file string, clPort string, srvPort string) (map[string]string, error) {
	// from https://stackoverflow.com/questions/24999079/reading-csv-file-in-go
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	ports := map[string]string{}

	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			if ports[clPort] != srvPort {
				PrintErr(errors.New("Local client and server ports don't match"))
			}
			delete(ports, clPort)
			return ports, err
		}

		ports[row[0]] = row[1]
	}
}

// ParseGroupMembersText parses the supplied text group member file
// https://stackoverflow.com/questions/36111777/how-to-read-a-text-file
func ParseGroupMembersText(file string, myPort string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	res, err := ioutil.ReadAll(f)
	ports := strings.Split(string(res), ",")
	// fmt.Println(ports)

	var j int
	flag := false
	for i, port := range ports {
		if port == myPort {
			flag = true
			j = i
			break
		}
	}

	// https://stackoverflow.com/questions/25025409/delete-element-in-a-slice
	if flag {
		ports[j] = ports[len(ports)-1] // Replace it with the last one. CAREFUL only works if you have enough elements.
		ports = ports[:len(ports)-1]   // Chop off the last one.
	}

	return ports, nil
}

// RPCClient makes an RPC connection
func RPCClient(channel chan *rpc.Client, logger *govec.GoLog, port string, who string) {
	options := govec.GetDefaultLogOptions()

	fmt.Println(who + "Starting RPC Connection to " + port)
	client, err := vrpc.RPCDial("tcp", "127.0.0.1:"+port, logger, options)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(who + "Connection made to " + port)
	channel <- client
}

// PrintMsg prints message to console from a replica
func PrintMsg(no int, msg string) {
	fmt.Println("REPLICA " + strconv.Itoa(no+1) + ": " + msg)
}

// PrintErr prints error
// from https://github.com/DistributedClocks/GoVector/blob/master/example/ClientServer/ClientServer.go
func PrintErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
