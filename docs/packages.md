# Package Specific Information
Jump to: [Client](#Client) | [Checker](#Checker) | [Crawler](#Crawler) | [Util](#Util) | [Docs](#Docs)

## CvRDT, CmRDTO, CmRDTC, Zero
Those folders contain the server code for the three implemernations of CRDTs studied in this project, as well as the "zero" implementation which makes no attemps to achieve any consistency.

* **cvrdt**: implementation of cvrdt
* **cmrdto**: implementation of cmrdt via standard casual broadcast (optimistic)
* **cmrdtc**: implementation of cmrdt via queueing operations (conservative)

### Starting the servers
* `ports.csv` must list all ips addresses, ports and database ports for all replicas in the group (one line per replica, in this order, separated by commas)
* then start any server by running `go run . [replicaNo] [communicationPort] [databasePort] [emulatedDelay] [verbose?]`
for example `go run . 1 8001 27017 100 2` 
* the `databasePort` and `emulatedDelay` parameters are useful when running several servers on a single machine, an `emulatedDelay` value of greater than zero will add random delays of the given value $\pm$ 20% in ms to all RPC calls
* if the verbose parameter is set to **1**, then the server will collect debugging information in logs, if it is set to **2** then the server will print information to console in addition to writing to logs, otherwise, it should be set to **0** when running performance evals)
* the CmRDT-C implementation takes the maximum queue length as an additional parameter for example `go run . 1 8001 27017 100 2 20` sets the maximum queue length to 20 OpNodes
* the CVRDT implementation takes a `y/n` if garbage collection should run as an additional parameter, for example `go run . 1 8001 27017 100 2 y` indicates the CvRDT will run with garbage collection

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
Implementation of client that sends commands to cvrdt/cmrdt servers along with various test sets of commands.

Start client by running `go run . [delayBetweenCommands] [timeSetting] [runMongoTest] [runRemoves] [terminateReplica]`
* where `delayBetweenCommands` is the time interval between commands send by the client (in ms)
* where `timeSetting` is the time interval between states exchanges by CvRDT replicas (in ms) or is the time interval between sending no-ops in the CmRDTC implementation
* where `runMongoTest` set to **y** indicates that the client will run tests that test MongoDb's native replication framework
* where `runRemoves` set to **y** indicates if the main test should test key/value removals as well (as opposed to just inserts)
* where `terminateReplica` set to **y** indicates if the main test should terminate replica once it is done (which runs lookup in the case of CmRDT-O and writes logs to files in all implementations)

Specific tests packages are implemented in `maintest.go`, `quicktest.go` and `wikitest.go`

## Checker
A program that checks consistency of replica's databases after a test run, downloads the database into CSV files labelled by replicas' numbers.

Start client by running `go run . [drop] [cvrdt/cmrdt]`
* if `drop` is equal to **y**, then tester will clear all databases to prepare the replicas for the following run
* if `cvrdt/cmrdt` is equal to **cv**, then tester will additionally download the positive and negative collections of the CvRDT servers and save those to CSV; otherwise, if it is equal to **cm** then tester will additionally download the CmRDT dynamic collection

## Crawler
A program that downloads Wikipedia pages to be used as test sets.

Start crawler by running `go run . [maxNoLinks] [maxDepth]`
* where `maxNoLinks` is the maximum number of outgoing links to follow from any given page
* where `maxDepth` is the maximum depth of the graph

## Util
Methods that are believed to be common to all implementations, perhaps more could be extracted into this package; however, it might start to obscure the code's readability.

* **util.go**: generic methods, such as parsing the group memberhship CSV file, finding max/mins, etc.
* **connops.go**: connection operations (e.g. connect one replica to another, connect a replica to a local database) that are common to all implementations
* **dbsops.go**: basic database operations that are common to all implementations

## Docs
MATLAB code that is used to generate figures, LaTeX report code, random notes, etc.
