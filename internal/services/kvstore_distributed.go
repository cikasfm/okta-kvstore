package services

import (
	"errors"
)

type KeyValueStore struct {
}

func NewKeyValueStore() (IKeyValueStore, error) {
	return &KeyValueStore{}, nil
}

func (s *KeyValueStore) Set(key string, value string) error {
	return errors.New("not implemented")
}

func (s *KeyValueStore) Get(key string) (string, error) {
	return "", errors.New("not implemented")
}

func (s *KeyValueStore) Delete(key string) error {
	return errors.New("not implemented")
}
