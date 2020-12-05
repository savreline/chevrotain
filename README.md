# Chevrotain
* **cvrdt**: implementation of cvrdt
* **cmrdtb**: implementation of cmrdt via standard casual broadcast
* **cmrdtq**: implementation of cmrdt via queueing operations
* **client**: implementation of client that sends commands to cvrdt/cmrdt servers along with various test sets of commands
* **tester**: program that checks consistency of replica's databases after a test run
* **crawler**: program that downloads Wikipedia pages to be used as test sets
* **util**: common methods, such as connecting a replica to the database, connecting replicas to each other, finding max/mins of values, etc.

start all servers by running `go run . [replicaNo] [communicationPort] [databasePort]` \
for example `go run . 1 8001 27017`

start client by running `go run . [delayBetweenCommands] [timeSetting]` \
where *delayBetweenCommands* is the time interval between commands send by the client (in ms)
where *timeSetting* is the time interval between states exchanges by CvRDT replicas (in ms)

`ports.csv` must list all ports and addresses for all replicas in the group as servers parse this file 

## cvrdt
* **server.go**: starts the server, contains `InitReplica` and `TerminateReplica` methods for `RPCExt`
* **dbops.go**: contains all other `RPCExt` methods and all operations with the local database
* **rpcint.go**:  
* **merges.go**: methods that
    1. merge incoming states and
    2. merge collections (i.e. positive and negative collection according to logical clocks)

## cmrdtb

## cmrdtq
