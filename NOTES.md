# Distributed Key-Value Store

## Descision

This project will try to implement and showcase development skills using GoLang

Author: [Zilvinas Vilutis](https://www.linkedin.com/in/zvilutis/)

## Descision & implementation process

### REST API & Scaffolding

First create a REST API endpoint using `gin-gonic` library and related tests.

The API endpoint impl [internal/api/kvstore_api.go](internal/api/kvstore_api.go) allows injection of a backing service
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

~~I'm leaning towards using [Zookeeper](https://zookeeper.apache.org/) for service discovery
while still researching the [raft](https://github.com/hashicorp/raft) framework~~

I'm not familiar with Raft framework but since it was suggested and I've discovered
that the example implementation contains the exact backing store required for this project
I will not hesitate and pull it in directly from: https://github.com/otoolep/hraftd/blob/master/store/store.go

### Deployment

I chose `Docker` platform for deployment lightweight, `alpine-linux` based `golang` containers
which using `docker compose` can be quickly built & deployed on any machine running `docker`.

In a docker environment internal network is created and so raft cluster endpoint is hidden
from outside of the container and exposed only to peer nodes, while REST API endpoints
are exposed outside.

Also, a docker volume is used to persist the `rockdb` data to file and reuse even after
restarting the container.

Kubernetes might've been a better option for easier scaling, but due to it's
design it's hiding any data about the internal nodes HTTP addresses, which is
crucial to create a raft cluster.

Kubernetes has advantages of scaling and load balancing the nodes and its
doing it's own service discovery, which is not what we need when we need
leader election and internal data replication.

In order to make this work with `kube` it might require some special handling
of routing `Set` and `Delete` requests, while `Get` could be run on any node
across the cluster. One way to get around this
( assuming there's a way to create a cluster ) is internally find the leader
in the raft cluster and forward the call.

That still requires service discovery, so for simplicity
I'm using `docker compose`.

### Project structure

There's no plan to extend this project and make parts of it usable & exposed
so most of its code is under `internal`.

All services provide a way to inject dependencies using a constructor, which
makes it easy to provide mock dependency implementations for unit testing.

## Summary

I'm using the backing store using raft from: https://github.com/otoolep/hraftd/blob/master/store/store.go that seems to be addressing raft cluster setup, electing a leader and also
snapshot replication, as soon as a new node joins the cluster it'll replicate the latest
snapshot. At this point the connection is insecure, so technically anything can connect to it,
which might cause some data leaks, but this was out of scope of the excercise.

On top of that I've implemented the REST API endpoint according to spec, covered majority of the edge cases.

Deployment to up to (but not limited to) 4 nodes using `docker compose`, or running manually by creating the raft cluster.

Storage either in memory using a simple hash map or on file, based on the parameters mentioned in [README.md](README.md)

When a node joins the cluster it'll replicate the latest snapshot from the leader. When a leader disconnects, a new leader is elected.

`Set` & `Delete` commands can only be executed on the leader, `Get` will work on all nodes.

Benchmark testing starts raft cluster ( single node ) so it needs ports to be opened and also it'll use quite some resources ( CPU & memory ) to run.

