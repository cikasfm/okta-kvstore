package services

import (
	"codesignal/internal/store"
	"errors"
	"fmt"
	"github.com/hashicorp/raft"
	"log"
	"os"
)

var (
	CodeRaftError     = "store error"
	CodeRaftNotLeader = "store not leader"
)

type RaftKeyValueStore struct {
	store  *store.Store
	logger *log.Logger
}

func NewRaftKeyValueStore(store *store.Store) (IKeyValueStore, error) {
	return &RaftKeyValueStore{
		store:  store,
		logger: log.New(os.Stderr, "[store-api] ", log.LstdFlags),
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

	existing, err := s.Get(key)
	if err != nil {
		s.logger.Println(fmt.Sprintf("error setting key to Raft store '%s': %s", key, err.Error()))
		return handleRaftError(err)
	}
	if existing != "" {
		return ServiceError{
			Code:    CodeKeyExists,
			Message: fmt.Sprintf("key %s already exists", key),
		}
	}

	// set key to Raft store
	err = s.store.Set(key, value)

	if err != nil {
		s.logger.Println(fmt.Sprintf("error setting key to Raft store '%s': %s", key, err.Error()))
		return handleRaftError(err)
	}

	return nil
}

func (s *RaftKeyValueStore) Get(key string) (string, error) {
	err := validateInput(key, KeyMaxLength)
	if err != nil {
		return "", err
	}

	value, err := s.store.Get(key)
	if err != nil {
		s.logger.Println(fmt.Sprintf("error getting key from Raft store %s: %s", key, err.Error()))
		return "", handleRaftError(err)
	}

	return value, nil
}

func (s *RaftKeyValueStore) Delete(key string) error {
	err := validateInput(key, KeyMaxLength)
	if err != nil {
		return err
	}

	err = s.store.Delete(key)
	if err != nil {
		s.logger.Println(fmt.Sprintf("error deleting key from Raft store %s: %s", key, err.Error()))
		return handleRaftError(err)
	}

	return nil
}

func handleRaftError(err error) error {
	if errors.Is(err, raft.ErrNotLeader) {
		return ServiceError{
			Code:    CodeRaftNotLeader,
			Message: "Unable to execute command: Node is not the leader",
		}
	}
	if errors.Is(err, raft.ErrNotVoter) {
		return ServiceError{
			Code:    CodeRaftNotLeader,
			Message: "Unable to execute command: Node is not a voter",
		}
	}
	if "not leader" == err.Error() {
		return ServiceError{
			Code:    CodeRaftNotLeader,
			Message: "Unable to execute command: Node is not the leader",
		}
	}
	return ServiceError{
		Code:    CodeRaftError,
		Message: "Unable to execute command: Internal Server Error",
	}
}
