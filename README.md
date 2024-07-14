# Distributed Key-Value Store

## Summary

This project will try to implement and showcase development skills using GoLang

Author: [Zilvinas Vilutis](https://www.linkedin.com/in/zvilutis/)

## Running

To run the program you can use `make` commands:

### Testing

#### Unit Tests

```shell
make test
```

### Running

Running a single instance.

```shell
make run
```

OR

```shell
./main.sh
```

This will run an instance listening on port `8080` and running `raft` on port `7000`.

To change this config - modify the environment variables in the [makefile](makefile)

**Config parameters:**

```shell
NODE_ID=app1 # <- node name
RAFT_BIND_ADDR=127.0.0.1:7000 # <- raft server bind address/port
RAFT_JOIN_ADDR=app1:8081 # <- URL of the Raft leader to join
RAFT_DIR=/usr/share/raft # <- Raft data storage folder. If set - file system will be used for storage
```

`NODE_ID` (required)
: node name to be used in raft cluster to identify nodes

`RAFT_BIND_ADDR` (optional, default=127.0.0.1:7000)
: raft server bind address/port for cluster operations

`RAFT_JOIN_ADDR` (optional, default=)
: URL of the Raft leader to join. When not set - the node will auto-elect as Leader. Otherwise - will try to join the cluster at the given address.

`RAFT_DIR` (optional, default=)
: Data storage folder. When not set - in memory storage used. Otherwise, `rockdb` files will be stored in this directory. 

#### Cluster

**Using Docker Compose**

The easiest way to run nodes in the cluster you can use the docker compose script:

```shell
make docker
```

Note: when peer nodes start while node1 has not yet selected itself as a leader, they might fail.

To avoid this, you can run 4 nodes one-by-one, starting with the first one:

```shell
# build & run as a daemon process
docker compose up app1 -d --build
sleep 2s
# wait a few seconds, then run as a daemon processes
docker compose up app2 app3 app4 -d
```

To stop `[all]` containers, run
```shell
docker compose down
```

**Running manually**

To start a cluster manually, first we need to build the go project:

```shell
make build
```

Then start nodes in separate processes:

```shell
# node1:
NODE_ID=app1 PORT=8081 RAFT_BIND_ADDR=127.0.0.1:7001 ./main
# node2 points to node1 leader:
NODE_ID=app2 PORT=8082 RAFT_BIND_ADDR=127.0.0.1:7002 RAFT_JOIN_ADDR=localhost:8081 ./main
# node3 points to node1 leader:
NODE_ID=app2 PORT=8083 RAFT_BIND_ADDR=127.0.0.1:7003 RAFT_JOIN_ADDR=localhost:8081 ./main
# node4 points to node1 leader:
NODE_ID=app2 PORT=8084 RAFT_BIND_ADDR=127.0.0.1:7004 RAFT_JOIN_ADDR=localhost:8081 ./main
```

#### Manual Tests

Run the following commands once deployed

**Set Key**:

```shell
curl -X POST 'http://localhost:8080/key' \
--header 'Content-Type: application/json' \
--data '{
    "key1": "value1"
}'
```

**Get Key**:

```shell
curl -X GET 'http://localhost:8080/key/key1'
```

**Get Key**:

```shell
curl -X DELETE 'http://localhost:8080/key/key1'
```