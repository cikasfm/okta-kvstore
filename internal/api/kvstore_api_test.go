package api

import (
	"codesignal/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestKeyValueStoreApi_deleteValue(t *testing.T) {
	inMemoryKeyValueStore := services.NewInMemoryKeyValueStore()

	_ = inMemoryKeyValueStore.Set("testKey", "testValue")

	tests := []struct {
		name       string
		key        string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "test delete existing value",
			key:        "testKey",
			wantStatus: 200,
			wantBody:   "{\"message\":\"key deleted successfully\"}",
		},
		{
			name:       "test delete non-existing value",
			key:        "nonExistent",
			wantStatus: 404,
			wantBody:   "{\"message\":\"key nonExistent not found\"}",
		},
		{
			name:       "test delete non-existing value",
			key:        "",
			wantStatus: 404,
			wantBody:   "404 page not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h := &KeyValueStoreApi{storage: inMemoryKeyValueStore}
			c, e := gin.CreateTestContext(w)

			req, _ := http.NewRequestWithContext(c, "DELETE", fmt.Sprintf("/key/%s", tt.key), nil)

			e.DELETE("/key/:key", h.DeleteValue)

			e.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestKeyValueStoreApi_getByKey(t *testing.T) {

	inMemoryKeyValueStore := services.NewInMemoryKeyValueStore()

	_ = inMemoryKeyValueStore.Set("testKey", "testValue")

	tests := []struct {
		name       string
		key        string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "test get by existing key",
			key:        "testKey",
			wantStatus: 200,
			wantBody:   "{\"testKey\":\"testValue\"}",
		},
		{
			name:       "test get by non-existing key",
			key:        "nonExistent",
			wantStatus: 404,
			wantBody:   "{\"message\":\"key nonExistent not found\"}",
		},
		{
			name:       "test get with empty key",
			key:        "",
			wantStatus: 404,
			wantBody:   "404 page not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h := &KeyValueStoreApi{storage: inMemoryKeyValueStore}
			c, e := gin.CreateTestContext(w)

			req, _ := http.NewRequestWithContext(c, "GET", fmt.Sprintf("/key/%s", tt.key), nil)

			e.GET("/key/:key", h.GetByKey)

			e.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestKeyValueStoreApi_setValue(t *testing.T) {

	inMemoryKeyValueStore := services.NewInMemoryKeyValueStore()

	_ = inMemoryKeyValueStore.Set("testKey", "testValue")

	tests := []struct {
		name       string
		body       io.Reader
		wantStatus int
		wantBody   string
	}{
		{
			name:       "test set by existing key",
			body:       strings.NewReader("{\"testKey\":\"testValue\"}"),
			wantStatus: 409,
			wantBody:   "{\"message\":\"key testKey already exists\"}",
		},
		{
			name:       "test set by non-existing key",
			body:       strings.NewReader("{\"newKey\":\"newValue\"}"),
			wantStatus: 200,
			wantBody:   "{\"newKey\":\"newValue\"}",
		},
		{
			name:       "test set with empty body",
			body:       nil,
			wantStatus: 400,
			wantBody:   "{\"message\":\"empty request body\"}",
		},
		{
			name:       "test set with not a json object",
			body:       strings.NewReader("not a json object"),
			wantStatus: 400,
			wantBody:   "{\"message\":\"request is not a valid json\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h := &KeyValueStoreApi{storage: inMemoryKeyValueStore}
			c, e := gin.CreateTestContext(w)

			req, _ := http.NewRequestWithContext(c, "POST", "/key", tt.body)

			e.POST("/key", h.SetValue)

			e.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestNewKeyValueStoreApi(t *testing.T) {
	type args struct {
		service services.IKeyValueStore
	}
	tests := []struct {
		name string
		args args
		want *KeyValueStoreApi
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewKeyValueStoreApi(tt.args.service), "NewKeyValueStoreApi(%v)", tt.args.service)
		})
	}
}

func Test_handleError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := handleError(tt.args.err)
			assert.Equalf(t, tt.want, got, "handleError(%v)", tt.args.err)
			assert.Equalf(t, tt.want1, got1, "handleError(%v)", tt.args.err)
		})
	}
}
