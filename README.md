# Chevrotain
Chevrotain is a replicated key value store that achieves eventual consistency through the use of a conflict-free replicated data types (CRDTs). This project implements and evaluates performances of three different design approaches to the implementation of Chevrotain. One of the approaches is based on a state-based CRDT model (CvRDT), while the other two approaches are based on an operation-based CRDT model (CmRDT), either with or without limited synchronization.

## CvRDT, CmRDTO, CmRDTC, Zero
Those folders contain the server code for the three implemernations of CRDTs studied in this project, as well as the "zero" implementation which makes no attemps to achieve any consistency

* **cvrdt**: implementation of cvrdt
* **cmrdto**: implementation of cmrdt via standard casual broadcast (optimistic)
* **cmrdtc**: implementation of cmrdt via queueing operations (conservative)

### Starting the servers
* `ports.csv` must list all ips addresses, ports and database ports for all replicas in the group (one line per replica, in this order, separated by commas)
* then start any server by running `go run . [replicaNo] [communicationPort] [databasePort] [emulatedDelay] [verbose?]`
for example `go run . 1 8001 27017 100 v` (the *databasePort* and *emulatedDelay* parameters are useful when running several servers on a single machine, an *emulatedDelay* value of greater than zero will add random delays of the given value +-20% in ms to all RPC calls; if the verbose parameter is set to "v" then the server will print to console and write to logs, otherwise, it should be set to "n" when running performance evals)
* the CVRDT implementation takes in an extra parameters, a y/n if garbage collection should run, for example `go run . 1 8001 27017 100 v y` indicates the CvRDT will run with garbage collection

### Specific files
* **server.go**: starts the server, initializes all variables and data structures; contains `InitReplica` and `TerminateReplica` methods for `RPCExt`
* **rpcext.go**: all other `RPCExt` methods, in particular, the InsertKey/InsertValue/RemoveKey/RemoveValue APIs
* **rpcint.go**: all `RPCInt` methods
* **dbops.go**: contains all methods that work with the local database
* **merges.go (CvRDT only)**: methods that
    1. merge incoming states and
    2. merge collections (i.e. positive and negative collection according to logical clocks)
* **queue.go (CmRDTC only)**: all methods that manage and process the OpNode queue, including insertion of incoming OpNode and processing of blocks of concurrent operations

In some implementations **rpcext.go** and **rpcint.go** are combined into **rpc.go**. Files **dbops_test.go** and **queue_test.go** contain some smoke tests to test correctness of the database and queueing operations in isolation.

## Client
implementation of client that sends commands to cvrdt/cmrdt servers along with various test sets of commands

start client by running `go run . [delayBetweenCommands] [timeSetting]` \
where *delayBetweenCommands* is the time interval between commands send by the client (in ms) \
where *timeSetting* is the time interval between states exchanges by CvRDT replicas (in ms) or is the time interval between sending no-ops in the CmRDTC implementation

specific tests packages are implemented in `test1.go`, `test2.go` and `wikitest.go`

## Tester
a program that checks consistency of replica's databases after a test run, downloads the database into CSV files labelled by replicas' numbers

start client by running `go run . [drop] [cvrdt/cmrdt]` \
if *drop* is equal to 1, then tester will clear all databases to prepare the replicas for the following run \
if *cvrdt/cmrdt* is equal to "cv", then tester will additionally download the positive and negative collections of the CvRDT servers and save those to CSV; otherwise, if it is equal to "cm" then tester will additionally download the CmRDT dynamic collection

## Crawler
a program that downloads Wikipedia pages to be used as test sets

start crawler by running `go run . [maxNoLinks] [maxDepth]` \
where *maxNoLinks* is the maximum number of outgoing links to follow from any given page \
where *maxDepth* is the maximum depth of the 

## Util
methods that are believed to be common to all implementations, perhaps more could be extracted into this package; however, it might start to obscure the code's readability

* **util.go**: generic methods, such as parsing the group memberhship CSV file, finding max/mins, etc.
* **connops.go**: connection operations (e.g. connect one replica to another, connect a replica to a local database) that are common to all implementations
* **dbsops.go**: basic database operations that are common to all implementations

## Docs
MATLAB code that is used to generate figures, LaTeX report code, random notes, etc.
