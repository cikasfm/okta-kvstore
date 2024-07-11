package services

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
)

type IKeyValueStore interface {
	//set the value to storage by key
	Set(key string, value string) error
	//get the value from storage by key
	Get(key string) (string, error)
	//delete the value from storage by key
	Delete(key string) error
}

type KeyValueStore struct {
}

func NewKeyValueStore() (IKeyValueStore, error) {
	return &KeyValueStore{}, nil
}

type ServiceError struct {
	Code    string
	Message string
}

func (e ServiceError) Error() string {
	return e.Message
}

var (
	CodeKeyNotFound = "key not found"
	CodeKeyExists   = "key already exists"
)

func (s *KeyValueStore) Set(key string, value string) error {
	return errors.New("not implemented")
}

func (s *KeyValueStore) Get(key string) (string, error) {
	return "", errors.New("not implemented")
}

func (s *KeyValueStore) Delete(key string) error {
	return errors.New("not implemented")
}

// in-memory map cache service
type InMemoryKeyValueStore struct {
	cache map[string]string
}

func NewInMemoryKeyValueStore() IKeyValueStore {
	return &InMemoryKeyValueStore{
		cache: make(map[string]string),
	}
}

func (s *InMemoryKeyValueStore) Set(key string, value string) error {
	if _, ok := s.cache[key]; ok {
		return ServiceError{
			Code:    CodeKeyExists,
			Message: spew.Sprintf("key %s already exists", key),
		}
	} else {
		s.cache[key] = value
		return nil
	}
}

func (s *InMemoryKeyValueStore) Get(key string) (string, error) {
	if value, ok := s.cache[key]; ok {
		return value, nil
	} else {
		return "", ServiceError{
			Code:    CodeKeyNotFound,
			Message: spew.Sprintf("key %s not found", key),
		}
	}
}

func (s *InMemoryKeyValueStore) Delete(key string) error {
	if _, ok := s.cache[key]; ok {
		delete(s.cache, key)
		return nil
	} else {
		return ServiceError{
			Code:    CodeKeyNotFound,
			Message: spew.Sprintf("key %s not found", key),
		}
	}
}
