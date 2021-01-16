# Chevrotain: a Replicated Key-Value Stores
Chevrotain is a replicated key value store that achieves eventual consistency through the use of a conflict-free replicated data types (CRDTs). This project implements and evaluates performances of three different design approaches to the implementation of Chevrotain. One of the approaches is based on a state-based CRDT model (CvRDT), while the other two approaches are based on an operation-based CRDT model (CmRDT), either with or without limited synchronization.

The author appreciates the insight and feedback received throughout this project from [Prof. Ivan Beschastnikh](https://www.cs.ubc.ca/~bestchai/) and ability to use Microsoft Azure Education credits.

**Main Golang Libraries:**
[net/rpc](https://golang.org/pkg/net/rpc/), 
[MongoDB](go.mongodb.org/mongo-driver/mongo),
[BSON](go.mongodb.org/mongo-driver/bson),
[GoVector](https://github.com/DistributedClocks/GoVector) \
**Other Libraries Used**:
[net/http](https://golang.org/pkg/net/http/), 
[net/html](golang.org/x/net/html),
[encoding/csv](https://golang.org/pkg/encoding/csv/),
[os/signal](https://golang.org/pkg/os/signal/),
[windows](https://pkg.go.dev/golang.org/x/sys/windows)

The key points of the [full fechnical report (PDF)](/docs/report/report.pdf) are summarized below.

---

Jump to [Background](#Background) | [Test Methodology](#Test-Methodology) 
| [Results](#Results)

## API
All communication between the client and any one of the replicas is done via the RPCExt object. The following methods can be called on this object. The parameters passed to the `InitReplica` method are implementation dependent. For example, in the CvRDT implementation, the `timeInt` parameter sets the time intervals between garbage collection, while the `bias` parameter is a struct that sets the user-defined bias in case of simultaneous InsertKey\RemoveKey and InsertValue\RemoveValue calls. All communication between the replicas is done via the RPCInt object, and the APIs there are implementation dependent.
* `InitReplica(timeInt int, bias Bias)`
* `InsertKey(key string)`
* `InsertValue(key string, value string)`
* `RemoveKey(key string)`
* `RemoveValue(key string, value string)`
* `TerminateReplica()`

## Background
**CvRDT**

[this](https://hal.inria.fr/inria-00555588/document)

<img src="docs/report/Fig3CvRDT2.png" width="400">

**CmRDT**

[this](https://hal.inria.fr/inria-00609399v1/document)

## Test Methodology

## Results

## Additional Information
[Package Specific Information](/docs/packages.md)
