package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func SetupStore() *Store {
	logger := log.New(os.Stderr, "[raft] ", log.LstdFlags)
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		logger.Fatal("NODE_ID environment variable is required")
	} else {
		logger.Printf("NODE_ID=%s \n", nodeID)
	}

	// static when running on docker
	raftBindAddr := os.Getenv("RAFT_BIND_ADDR")
	if raftBindAddr == "" {
		raftBindAddr = "127.0.0.1:7000"
	}
	logger.Printf("RAFT_BIND_ADDR=%s \n", raftBindAddr)

	raftAddr, err := net.ResolveTCPAddr("tcp", raftBindAddr)
	if err != nil {
		logger.Fatalf("Failed to resolve TCP raftAddr: %v", err)
	} else {
		logger.Printf("raftAddr=%v \n", raftAddr)
	}

	raftDir := os.Getenv("RAFT_DIR")

	inMemory := raftDir == ""

	var store *Store

	if inMemory {
		store = NewStore(true)
	} else {
		store = NewStore(false)

		path := filepath.Join(raftDir, nodeID)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			logger.Println(fmt.Sprintf("Creating directory: %s", path))
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				logger.Fatalf("Failed to create dir %s: %v", path, err)
			}
		} else {
			logger.Printf("path=%v \n", path)
		}

		store.RaftDir = path
	}

	joinAddr := os.Getenv("RAFT_JOIN_ADDR")

	store.RaftBind = raftAddr.String()

	err = store.Open(joinAddr == "", nodeID)
	if err != nil {
		logger.Fatalf("Failed to start raft: %v", err)
	}

	// If join was specified, make the join request.
	if joinAddr != "" {
		if err := join(joinAddr, raftAddr.String(), nodeID); err != nil {
			log.Fatalf("failed to join node at %s: %s", joinAddr, err.Error())
		}
	}

	logger.Println(fmt.Sprintf("Raft started successfully on raftAddr %s", raftAddr))

	return store
}

func join(joinAddr, raftAddr, nodeID string) error {
	b, err := json.Marshal(map[string]string{"address": raftAddr, "node_id": nodeID})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/join", joinAddr), "application-type/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
