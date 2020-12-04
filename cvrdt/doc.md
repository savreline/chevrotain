# CvRDT
#### server.go
* **main:** 
    1. parse command line arguments, parse group member info, initialize data structures
    2. connect to mongo, pre-allocate keys document, init RPC, start background processes
* **InitReplica:**
    1. set passed in settings
    2. activate background processes
    3. make RPC connections to other replicas
* **TerminateReplica:**
    1. write logs to text file

#### dbops.go
* **RPCExt**:
    1. InsertLocalRecord
    2. emulateDelay
* **InsertLocalRecord**:
    1. tick the clock
    2. if no record is supplied, make one
    otherwise, check if record already been inserted
    3. if no key document exists, make one
    4. push the record
    5. print to console

#### rpcint.go
* **StateArgs**:
	1. PosState, NegState []util.CvDoc
	2. SrcPid
	3. Timestamp
* **MergeState**:
    1. merge clock
    2. merge each collection
* **runSE/runGC**
* **broadcast**
    1. tick the clock
    2. download state
    3. broadcast
* **gc**

#### merges.go
* **mergeState**
    1. merge a collection: iterate over docs and record, call InsertLocalRecord
* **mergeCollections**
Note: there is a remote possibility of duplicate entries during merge state, those will be ignored during merge collections

#### client.go
1. parse command line arguments: delay between sending commands, replica settings time interval
2. parse group member info
3. run tests
    * connect client to replica and initialize replica
    * init map of latencies and wait group
    * send the commands
    * terminate replica
    * process collected peformance data

#### tester.go
1. parse command line arguments and group membership info
2. init data structures: pointers to all collections to fetch
3. connect to databases
4. download each collection
5. save each collection to CSV
