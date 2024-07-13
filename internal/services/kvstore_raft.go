package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/raft"
	"log"
	"time"
)

var (
	CodeRaftError = "raft error"
)

type RaftKeyValueStore struct {
	ra *raft.Raft
}

func NewRaftKeyValueStore(ra *raft.Raft) (IKeyValueStore, error) {
	return &RaftKeyValueStore{
		ra: ra,
	}, nil
}

func (s *RaftKeyValueStore) Set(key string, value string) error {
	err := validateInput(key, KeyMaxLength)
	if err != nil {
		return err
	}
	err = validateInput(value, ValueMaxLength)
	if err != nil {
		return err
	}
	data, err := json.Marshal(map[string]string{key: value})
	if err != nil {
		return err
	}

	applyFuture := s.ra.Apply(data, 10*time.Second)

	if err := applyFuture.Error(); err != nil {
		log.Println(fmt.Errorf("set error: %v", err))
		return ServiceError{
			Code:    CodeRaftError,
			Message: "Failed to apply command",
		}
	}
	return nil
}

func (s *RaftKeyValueStore) Get(key string) (string, error) {
	err := validateInput(key, KeyMaxLength)
	if err != nil {
		return "", err
	}

	// Raft does not support direct read-only queries
	f := s.ra.Apply([]byte(`{"operation":"get", "key":"`+key+`"}`), 10*time.Second)
	if err := f.Error(); err != nil {
		log.Println(fmt.Errorf("get error: %v", err))
		return "", errors.New("failed to apply command")
	}

	result := f.Response()
	value, ok := result.(string)
	if !ok {
		log.Println(fmt.Errorf("get error: %v", err))
		return "", errors.New("invalid raft response")
	}
	return value, nil
}

func (s *RaftKeyValueStore) Delete(key string) error {
	err := validateInput(key, KeyMaxLength)
	if err != nil {
		return err
	}
	return errors.New("not implemented")
}
