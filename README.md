# Chevrotain: a CRDT-Based Replicated Key-Value Store
Chevrotain is a replicated multi-primary key value store that achieves eventual consistency through the use of a conflict-free replicated data types (CRDTs). This project implements and evaluates performances of three different design approaches to the implementation of such a store. The three approaches are:
1. State-based CRDT model (CvRDT) with or without garbage collection
2. Operation-based CRDT model (CmRDT) without any synchronization
3. Operation-based CRDT model (CmRDT) with limited synchronization

Performance was evaluated by subjecting the three implementations to various loads while deployed on ten geographically distributed Azure VMs. Performance was also compared against MongoDB's built-in [replication service](https://docs.mongodb.com/manual/replication/) which follows the primary-backup model. The key points of the [full fechnical report (PDF)](docs/report/report.pdf) are summarized below.

This was an individual project completed in the fall of 2020 for the *Distributed Systems Abstractions* graduate course (CPSC 538B) at UBC. The insight and feedback received throughout this project from [Prof. Ivan Beschastnikh](https://www.cs.ubc.ca/~bestchai/) and ability to use Microsoft Azure Education credits is much appreciated.

**Main Golang Libraries:**
[net/rpc](https://golang.org/pkg/net/rpc/), 
[MongoDB](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo),
[BSON](https://pkg.go.dev/go.mongodb.org/mongo-driver/bson),
[GoVector](https://github.com/DistributedClocks/GoVector) \
**Other Libraries Used**:
[net/http](https://golang.org/pkg/net/http/), 
[net/html](https://pkg.go.dev/golang.org/x/net/html),
[encoding/csv](https://golang.org/pkg/encoding/csv/),
[os/signal](https://golang.org/pkg/os/signal/),
[windows](https://pkg.go.dev/golang.org/x/sys/windows)

---

Jump to [Background](#Background) | [Test Methodology](#Test-Methodology) 
| [Results](#Results)

## API
All communication between the client and any one of the replicas is done via the RPCExt object. The following methods can be called on this object. The parameters passed to the `InitReplica` method are implementation dependent. For example, in the CvRDT implementation, the `timeInt` parameter sets the time intervals between state exchanges, while the `bias` parameter is a struct that sets the user-defined bias in case of simultaneous InsertKey\RemoveKey and InsertValue\RemoveValue calls. All communication between the replicas is done via the RPCInt object, and the APIs there are implementation dependent.
* `InitReplica(timeInt int, bias Bias)`
* `InsertKey(key string)`
* `InsertValue(key string, value string)`
* `RemoveKey(key string)`
* `RemoveValue(key string, value string)`
* `TerminateReplica()`

## Background
**CvRDT**

The implementation largely follows the approach described in section 3.3.3 of [this](https://hal.inria.fr/inria-00555588/document) paper by Marc Shapiro et. all. 

<img src="docs/cvrdt.jpg" width="600">

**CmRDT**

The implementation largely follows the approach described in section 5 and figure 3 of [this](https://hal.inria.fr/inria-00609399v1/document) paper by Marc Shapiro et. all.

## Test Methodology
All implementations of the project were deployed on up to ten Azure D4s v3 VMs located in Canada Central, UK South, Japan East, Australia East and Brazil South zones. All implementations were subjected to a standard test that evenly distributed 1050 API calls between the given set of replicas. The rate at which the API calls were delivered to the replicas was varied and resulting end-to-end latency, consistency and time to reach steady state (CvRDT only) were measured. In a separate experiment, MongoDB's built-in replication service was set-up between the same replicas and the primary replica was subjected to the same API calls.

## Results
* 

## Additional Information
[Package Specific Information](docs/packages.md)
