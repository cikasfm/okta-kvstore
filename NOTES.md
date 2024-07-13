# Distributed Key-Value Store

## Summary

This project will try to implement and showcase development skills using GoLang

Author: [Zilvinas Vilutis](https://www.linkedin.com/in/zvilutis/)

## Descision & implementation process

### REST API & Scaffolding

First create a REST API endpoint using `gin-gonic` library and related tests.

The API endpoint impl [internal/api/kvstoreapi.go](internal/api/kvstoreapi.go) allows injection of a backing service
which in term allows independent testing of the API Endpoint without the implementation being completed.

## Distributed Key-Value Store design

### Storage

For storing key-value pairs we can use ( at least ) these options:

* In memory
  * Hash Map for random unstructured keys
  * Trie ( for Binary Search ) when keys can have common prefixes and so can be grouped by prefix.
  * Cache implementation ( e.g. Caffeine-like [goburrow/cache](https://github.com/goburrow/cache) )
    similar to a HashMap but optimized for KV cache usage
* External
  * File Storage ( local disk, external disk, network storage or even remote like s3 )
  * Database ( NoSQL ) or External Storage, like a Search Engine, Redis Cache
  * 3rd party service

And many others

Each has their cons and pros and can be a good fit based on the type of data.

In memory solutions are the fastest, but also have limitations:
- limited size
- loss of data ( in case of restart )

External storage on the other hand has low risk data storage and allows storage with no loss,
but is also potentially much slower.

### Replication

This is the interesting part.

There are multiple things to take into consideration:

* Replication type
  * full ( i.e. replicate full files / databases / storage )
    * on every write, OR
    * as a scheduled job
  * incremental
    * replicate on write, OR
    * replicate on read ( check neighbours for value, take the lastest one and update local version )
  * batch
    * collect batch updates and send to all neighbours based on a schedule and/or max batch size
* Locking
  * It's important to understand whether locking is needed for writing during a race condition:
    * When more than one node is trying to call `Set` or `Delete` at the same time 
      which would result in modification of a value
    * One option to solve this is by creating a `lock` on a row by key on all/consensus nodes prior to update,
      then update all nodes and `unlock`. Any other node trying to modify a value ( Set or Delete )
      would be blocked by the lock and should either wait for `unlock` or fail
  * In case of full replication - then full node data would have to be locked and all modify operations would be blocked
    * That can create a high consistency level but with a large penalty of performance

I'm leaning towards using [Zookeeper](https://zookeeper.apache.org/) for service discovery
while still researching the [raft](https://github.com/hashicorp/raft) framework 