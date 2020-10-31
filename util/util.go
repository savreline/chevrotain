package util

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

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
func ParseGroupMembersText(file string, port string) ([]string, error) {
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
	return strings.Split(string(res), ","), nil
}

// PrintErr prints error
// from https://github.com/DistributedClocks/GoVector/blob/master/example/ClientServer/ClientServer.go
func PrintErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
