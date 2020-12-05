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
    1. download both collections
    2. iterate over the documents in the positive collection
        * grab a pointer to the document in the negative collection
        * iterate over the records in the document
            * find all instances in positive and negative collections
            and determine the respective max timestamps
            * based on timestamps and settings bias, determine if the element
            need to be inserted or removed
            * remove all instances of the element from either collection
            * remove all instances of the element from the positive set iterating over
            * insert/remove the element into permament collection as need be
    3. iterate over the remaining documents in the negative collection
        * remove those key/values from the permament collection

    Notes: 
    1. there is a remote possibility of multiple documents with the same key, code can be adjusted to
    handle those when merging collections
    2. if remove wins, remove the element in all cases

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
