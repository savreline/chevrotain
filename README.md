# Chevrotain: a CRDT-Based Replicated Key-Value Store
Chevrotain is a replicated multi-primary key value store that achieves eventual consistency through the use of conflict-free replicated data types (CRDTs). This project implements and evaluates performances of three different design approaches to the implementation of such a store. The three approaches are:
1. State-based CRDT model (CvRDT) with or without garbage collection
2. Operation-based CRDT model (CmRDT-O) without any synchronization
3. Operation-based CRDT model (CmRDT-C) with limited synchronization

Performance was evaluated by subjecting the three implementations to various loads while deployed on ten geographically distributed Azure VMs. Performance was also compared against MongoDB's built-in [replication service](https://docs.mongodb.com/manual/replication/) which follows the primary-backup model. The key points of the [full technical report (PDF)](docs/report/report.pdf) are summarized below.

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
[os/signal](https://golang.org/pkg/os/signal/)

[![Go Report Card](https://goreportcard.com/badge/github.com/savreline/chevrotain)](https://goreportcard.com/report/github.com/savreline/chevrotain)

---

Jump to [Background](#Background) | [Test Methodology](#Test-Methodology) 
| [Results for a 3-Replica System](#3-replica-system) | [Scalability Results](#scalability-results)

## API
All communication between the client and any one of the replicas is done via the RPCExt object. The following methods can be called on this object. The parameters passed to the `InitReplica` method are implementation dependent. For example, in the CvRDT implementation, the `timeInt` parameter sets the time interval between state exchanges, while the `bias` parameter is a struct that sets the user-defined bias in case of concurrent InsertKey\RemoveKey and InsertValue\RemoveValue calls. All communication between the replicas is done via the RPCInt object, and the APIs for that object are implementation dependent.
* `InitReplica(timeInt int, bias Bias)`
* `InsertKey(key string)`
* `InsertValue(key string, value string)`
* `RemoveKey(key string)`
* `RemoveValue(key string, value string)`
* `TerminateReplica()`

## Background
There are two main approaches to maintaining consistency in replicated systems. The first approach is pessimistic. Some particular strategies are:
* lock down the entire system while a user makes changes to some replica and then propagate the changes to all other replicas
* allow multiple users to make changes to multiple replicas at the same time, but require the replicas to come to a universal consensus on the order and nature of changes
Both of those strategies are inefficient and degrade the system's availability. However, in both cases, perfect consistency is maintained at all times.

The second approach is optimistic and always allows multiple users to make changes to multiple replicas at the same time. Should the changes lead to conflicts, the conflicts are either resolved or the changes are rolled back. Optimistic systems are highly available but maintain only eventual consistency (as in, the states of replicas may temporarily diverge but always converge to an identical state later in time).

CRDTs is a technique of achieving eventual consistency by working with data in a way that avoids conflicts all together or resolves conflicts automatically. There are two approaches to implementing CRDTs: state-based (CvRDT) and operation-based (CmRDT). In the state-based approach, states of the replicas are maintained in a way that they could be merged in a conflict-free way at any point in time. In the operation-based approach, any updates done to the states of the replicas are done in a way to prevent conflicts at any point in time.

#### CvRDT
The implementation largely follows the approach described in section 3.3.3 of [this](https://hal.inria.fr/inria-00555588/document) paper by Marc Shapiro et. all. 

A key-value store could be represented by a set of keys and a set of values for each key. The state based approach maintains two sets for each set of elements: a "positive" set and a "negative" set. Each element in each set is tagged with a logical timestamp representing the relative time at which the element was added to the set. 

The states of two replicas are merged by taking the union of the respective pairs of sets. The positive and negative sets are merged by inserting an element into the merged set if and only if it is found in the positive set with a later timestamp than in the negative set. Should the timestamps of the element in the positive and negative sets be identical, a preset user-defined bias is used to tilt the result towards the presence or absence of the element in the merged set.

States are exchanged between replicas and are merged at preset intervals of time. Furthermore, to improve performance, elements of the positive and negative sets whose timestamps are below a "safe" timestamp are moved to the merged set (see figure 1). This process is known as "garbage collection". The "safe" timestamp is the minimum of the logical timestamps at all replicas at the given point in time.

#### Figure 1: Garbage collection in CvRDT
<img src="docs/cvrdt.jpg" width="600">

#### CmRDT-O
The implementation largely follows the approach described in section 5 and figure 3 of [this](https://hal.inria.fr/inria-00609399v1/document) paper by Marc Shapiro et. all and doesn't involve any synchronization.

An update taking place at one replica is propagated to all other replicas using a casual broadcast communication protocol (CBCAST). In particular, all operations are tagged with a vector clock and are processed in the order of their timestamps as they arrive at each replica. 

Updates are split into prepare-update and effect-update methods. The prepare-update method is side-effect free and snapshots the changes to be made to the data. This method takes place only at the replica at which the operation was initially delivered to and the noted changes are immediately applied to that replica by executing the effect-update method there. Once the operation has been delivered to all other replicas by CBCAST, the effect-update method is executed on those replicas as well. 

Some examples:
* For the insert key operation, the key is tagged with a unique id generated at the initiating replica and is inserted with that id into the databases at all replicas. 
* In the remove key operation, the prepare-update method gathers the unique ids of all instances of the key to be removed and those instances are removed at all replicas. This way, should an insert key operation occur while the key is being removed, the new instance of the key will not be removed.
* Should there be an insert value operation on a non-existing key, the value is inserted anyway into the internal database. Should the corresponding insert key operation never arrive, the value is hidden from any actual database queries.

#### CmRDT-C
In this implementation, all operations are tagged with a vector clock and are best-effort sorted according to those vector clocks on a queue maintained at each replica. Concurrent adds and removes are resolved as per used defined biases.

## Test Methodology
All implementations of the project were deployed on up to ten Azure D4s v3 VMs located in Canada Central, UK South, Japan East, Australia East and Brazil South zones. All implementations were subjected to a standard test that evenly distributed 1050 API calls between the given set of replicas. The rate at which the API calls were delivered to the replicas varied from 10 ops/s to 10 000 ops/s and resulting end-to-end latency, consistency and time to reach steady state (CvRDT only) were measured. In a separate experiment, MongoDB's built-in replication service was set-up between the same replicas and the primary replica was subjected to the same API calls.

## Results
### 3-Replica System
* CvRDT implementation with garbage collection (CvRDT-GC in the figure) performed best, maintaining latency of about 100ms under all loads.
* CvRDT implementation without garbage collection (CvRDT in the figure) performed second best, maintaining latency of about 100ms for throughput of up to 250 ops/s. Latency increased to about 200ms under a throughput of 10 000 ops/s.
* MongoDB's built-in replication service maintained latency less than 150ms for throughput of up to 100 ops/s. Latency increased to about 1s at throughput of 1000 ops/s. Throughput saturated at that point as well.
* CmRDT implementations performed worst, with CmRDT-O implementation demonstrating acceptable latency only for throughput less than or equal to 100 ops/s <sup>1</sup>. Performance of the CmRDT-C implementation was even less notable.
  
<sup>1</sup> 325ms at 100 ops/s, 765ms at 175 ops/s, 2.4s at 250 ops/s

#### Figure 2: Latency as a function of throughput for a 3-replica system
<img src="docs/report/Fig7TPLat3.png" width="500">

However, end-to-end latency measurements for CvRDT do not include time delays on the order of tens of seconds to exchange and merge states. Those delays are a function of the pre-set time interval that determines the frequency of state exchanges. In the CvRDT-GC implementation, at a throughput of 10 000 ops/s, delay decreases to 7s when state exchanges run every 100ms and increases to 22s when state exchanges run every 5000ms. Therefore, in a way, CvRDT implementation shifts the delays from the client to the server.

#### Figure 3: Merge time as a function of time between state exchanges in the CvRDT system
<img src="docs/report/Fig9CvRDTMerge3.png" width="500">

### Scalability Results
* CvRDT-GC implementation scaled to 10-replicas without any significant loss in performance.
* MongoDB's built-in replication service saturated at similar throughput thresholds (1000 ops/s in a 5-replica system and 750 ops/s in 7 and 10-replica systems).
* CmRDT-O implementation demonstrated unacceptable latency at slightly lower throughput thresholds (100 ops/s in a 5-replica system, 75 ops/s in a 7-replica system 50 ops/s in a 10-replica system)

#### Figure 4: Scalability of CvRDT-GC
<img src="docs/report/Fig10SCvRDTGC.png" width="500">

#### Figure 5: Scalability of MongoDB's Built-in Replication
<img src="docs/report/Fig13SMongo.png" width="500">

#### Figure 6: Scalability of CmRDT-O
<img src="docs/report/Fig12SCmRDTO.png" width="500">

## Additional Information
[Package Specific Information](docs/packages.md)
