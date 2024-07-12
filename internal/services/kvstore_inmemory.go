package services

import "fmt"

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
	err := validateInput(key, KeyMaxLength)
	if err != nil {
		return err
	}
	err = validateInput(value, ValueMaxLength)
	if err != nil {
		return err
	}
	if _, ok := s.cache[key]; ok {
		return ServiceError{
			Code:    CodeKeyExists,
			Message: fmt.Sprintf("key %s already exists", key),
		}
	} else {
		s.cache[key] = value
		return nil
	}
}

func (s *InMemoryKeyValueStore) Get(key string) (string, error) {
	err := validateInput(key, KeyMaxLength)
	if err != nil {
		return "", err
	}
	if value, ok := s.cache[key]; ok {
		return value, nil
	} else {
		return "", ServiceError{
			Code:    CodeKeyNotFound,
			Message: fmt.Sprintf("key %s not found", key),
		}
	}
}

func (s *InMemoryKeyValueStore) Delete(key string) error {
	err := validateInput(key, KeyMaxLength)
	if err != nil {
		return err
	}
	if _, ok := s.cache[key]; ok {
		delete(s.cache, key)
		return nil
	} else {
		return ServiceError{
			Code:    CodeKeyNotFound,
			Message: fmt.Sprintf("key %s not found", key),
		}
	}
}
