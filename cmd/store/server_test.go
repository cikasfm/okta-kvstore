package main

import (
	"codesignal/cmd/store/api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_setupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		method      string
		uri         string
		wantStatus  int
		wantContent string
	}{
		{
			name:        "Test setup router",
			method:      http.MethodGet,
			uri:         "/health",
			wantStatus:  http.StatusOK,
			wantContent: "\"OK\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setupRouter()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.uri, nil)
			got.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantContent, w.Body.String())
		})
	}
}

func Test_setupRoutes(t *testing.T) {
	type args struct {
		router  *gin.Engine
		handler *api.KeyValueStoreApi
	}
	tests := []struct {
		name string
		args args
		want *gin.Engine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, setupRoutes(tt.args.router, tt.args.handler), "setupRoutes(%v, %v)", tt.args.router, tt.args.handler)
		})
	}
}

func Test_startServer(t *testing.T) {
	type args struct {
		handler *api.KeyValueStoreApi
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startServer(tt.args.handler)
		})
	}
}
