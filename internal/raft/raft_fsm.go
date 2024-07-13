package raft

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/raft"
	"io"
	"log"
)

type fsm struct {
	store map[string]string
}

func (f *fsm) Apply(raftLog *raft.Log) interface{} {
	var command map[string]string
	if err := json.Unmarshal(raftLog.Data, &command); err != nil {
		log.Println(fmt.Errorf("failed to unmarshal log data: %v", err))
	}

	for key, value := range command {
		f.store[key] = value
	}
	return nil
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return &fsmSnapshot{store: f.store}, nil
}

func (f *fsm) Restore(rc io.ReadCloser) error {
	var store map[string]string
	if err := json.NewDecoder(rc).Decode(&store); err != nil {
		return err
	}
	f.store = store
	return nil
}

type fsmSnapshot struct {
	store map[string]string
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		data, err := json.Marshal(f.store)
		if err != nil {
			return err
		}
		if _, err := sink.Write(data); err != nil {
			return err
		}
		if err := sink.Close(); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		sink.Cancel()
	}
	return err
}

func (f *fsmSnapshot) Release() {}
