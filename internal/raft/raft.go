package raft

import (
	"fmt"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	raftDir = "./raft"
)

func SetupRaft() *raft.Raft {
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		log.Fatal("[raft] NODE_ID environment variable is required")
	} else {
		log.Printf("[raft] NODE_ID=%s \n", nodeID)
	}

	// static when running on docker
	raftBindAddr := os.Getenv("RAFT_BIND_ADDR")
	if raftBindAddr == "" {
		raftBindAddr = "127.0.0.1:7000"
	}
	log.Printf("[raft] RAFT_BIND_ADDR=%s \n", raftBindAddr)

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	address, err := net.ResolveTCPAddr("tcp", raftBindAddr)
	if err != nil {
		log.Fatalf("Failed to resolve TCP address: %v", err)
	} else {
		log.Printf("[raft] address=%v \n", address)
	}

	log.Println("Setting up raft transport")

	transport, err := raft.NewTCPTransport(raftBindAddr, address, 3, 10*time.Second, os.Stderr)
	if err != nil {
		log.Fatalf("Failed to create TCP transport for address %s: %v", raftBindAddr, err)
	} else {
		log.Printf("[raft] transport=%v \n", transport.LocalAddr())
	}

	path := filepath.Join(raftDir, nodeID)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		log.Println(fmt.Sprintf("Creating directory: %s", path))
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create dir %s: %v", path, err)
		}
	} else {
		log.Printf("[raft] path=%v \n", path)
	}

	// Setup Raft Stable storage
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, nodeID, "stableStore.db"))
	if err != nil {
		log.Fatalf("Failed to create stableStore: %v", err)
	} else {
		log.Println("[raft] stable store created")
	}

	// Setup Raft Log storage
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, nodeID, "logStore.db"))
	if err != nil {
		log.Fatalf("Failed to create logStore: %v", err)
	} else {
		log.Println("[raft] log store created")
	}

	snapshotStore, err := raft.NewFileSnapshotStore(filepath.Join(raftDir, nodeID), 1, os.Stderr)
	if err != nil {
		log.Fatalf("Failed to create snapshot store: %v", err)
	} else {
		log.Println("[raft] snapshot store created")
	}

	fsm := &fsm{store: make(map[string]string)}

	ra, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		log.Fatalf("Failed to create Raft: %v", err)
	} else {
		log.Println("[raft] instance created")
	}

	initCluster := os.Getenv("INIT_CLUSTER") == "true"
	initialServers := os.Getenv("INITIAL_SERVERS")

	if initCluster {
		log.Println("Initializing cluster")
		initialConfig := raft.Configuration{}
		for _, server := range strings.Split(initialServers, ",") {
			parts := strings.Split(server, ":")
			if len(parts) != 3 {
				log.Fatalf("Invalid server format in INITIAL_SERVERS: %s", server)
			}
			id := raft.ServerID(parts[0])
			address := raft.ServerAddress(parts[1] + ":" + parts[2])
			initialConfig.Servers = append(initialConfig.Servers, raft.Server{
				Suffrage: raft.Voter,
				ID:       id,
				Address:  address,
			})
		}

		f := ra.BootstrapCluster(initialConfig)
		if f.Error() != nil {
			log.Fatalf("Failed to bootstrap cluster: %v", f.Error())
		}
	}

	log.Println(fmt.Sprintf("Raft started successfully on address %s", address))

	return ra
}
