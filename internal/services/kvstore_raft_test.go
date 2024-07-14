package services

import (
	"codesignal/internal/store"
	"fmt"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"testing"
	"time"
)

func BenchmarkRaftKeyValueStore_Set(b *testing.B) {

	kvstore := setupStore(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		err := kvstore.Set(key, value)
		assert.NoError(b, err)
	}
}

func BenchmarkRaftKeyValueStore_Get(b *testing.B) {

	kvstore := setupStore(b)

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		err := kvstore.Set(key, value)
		assert.NoError(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		actual, err := kvstore.Get(key)
		assert.NoError(b, err)
		assert.Equalf(b, actual, value, "value should be equal")
	}
}

func BenchmarkRaftKeyValueStore_Delete(b *testing.B) {

	kvstore := setupStore(b)

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		err := kvstore.Set(key, value)
		assert.NoError(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		err := kvstore.Delete(key)
		assert.NoError(b, err)
	}
}

func setupStore(b *testing.B) IKeyValueStore {

	b.StopTimer()

	defer b.StartTimer()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf("localhost:%d", listener.Addr().(*net.TCPAddr).Port)
	fmt.Println("Using address:", address)

	defer listener.Close()

	s := store.NewStore(true)
	s.RaftBind = address
	err = s.Open(true, "test")
	assert.NoError(b, err)

	startTime := time.Now()
	for s.State() != raft.Leader.String() {
		time.Sleep(10 * time.Millisecond)
		if time.Since(startTime) > 2*time.Second {
			break
		}
	}

	kvstore, _ := NewRaftKeyValueStore(s)

	return kvstore
}

func TestNewRaftKeyValueStore(t *testing.T) {
	type args struct {
		store *store.Store
	}
	tests := []struct {
		name    string
		args    args
		want    IKeyValueStore
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRaftKeyValueStore(tt.args.store)
			if !tt.wantErr(t, err, fmt.Sprintf("NewRaftKeyValueStore(%v)", tt.args.store)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewRaftKeyValueStore(%v)", tt.args.store)
		})
	}
}

func TestRaftKeyValueStore_Delete(t *testing.T) {
	type fields struct {
		store  *store.Store
		logger *log.Logger
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RaftKeyValueStore{
				store:  tt.fields.store,
				logger: tt.fields.logger,
			}
			tt.wantErr(t, s.Delete(tt.args.key), fmt.Sprintf("Delete(%v)", tt.args.key))
		})
	}
}

func TestRaftKeyValueStore_Get(t *testing.T) {
	type fields struct {
		store  *store.Store
		logger *log.Logger
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RaftKeyValueStore{
				store:  tt.fields.store,
				logger: tt.fields.logger,
			}
			got, err := s.Get(tt.args.key)
			if !tt.wantErr(t, err, fmt.Sprintf("Get(%v)", tt.args.key)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Get(%v)", tt.args.key)
		})
	}
}

func TestRaftKeyValueStore_Set(t *testing.T) {
	type fields struct {
		store  *store.Store
		logger *log.Logger
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RaftKeyValueStore{
				store:  tt.fields.store,
				logger: tt.fields.logger,
			}
			tt.wantErr(t, s.Set(tt.args.key, tt.args.value), fmt.Sprintf("Set(%v, %v)", tt.args.key, tt.args.value))
		})
	}
}

func Test_handleRaftError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, handleRaftError(tt.args.err), fmt.Sprintf("handleRaftError(%v)", tt.args.err))
		})
	}
}
