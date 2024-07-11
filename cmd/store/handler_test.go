package main

import (
	"codesignal/cmd/store/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpHandler_deleteValue(t *testing.T) {
	w := httptest.NewRecorder()
	inMemoryKeyValueStore := services.NewInMemoryKeyValueStore()

	_ = inMemoryKeyValueStore.Set("testKey", "testValue")

	tests := []struct {
		name       string
		key        string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "deleteValue",
			key:        "testKey",
			wantStatus: 200,
			wantBody:   "{\"message\":\"success\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpHandler{storage: inMemoryKeyValueStore}
			c, e := gin.CreateTestContext(w)

			req, _ := http.NewRequestWithContext(c, "DELETE", "/key/testKey", nil)

			e.DELETE("/key/:key", h.deleteValue)

			e.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHttpHandler_getByKey(t *testing.T) {
	type fields struct {
		service services.IKeyValueStore
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpHandler{
				storage: tt.fields.service,
			}
			h.getByKey(tt.args.c)
		})
	}
}

func TestHttpHandler_setValue(t *testing.T) {
	type fields struct {
		service services.IKeyValueStore
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpHandler{
				storage: tt.fields.service,
			}
			h.setValue(tt.args.c)
		})
	}
}

func TestNewHttpHandler(t *testing.T) {
	type args struct {
		service services.IKeyValueStore
	}
	tests := []struct {
		name string
		args args
		want *HttpHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewHttpHandler(tt.args.service), "NewHttpHandler(%v)", tt.args.service)
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
