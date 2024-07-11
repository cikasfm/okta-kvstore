package main

import (
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
