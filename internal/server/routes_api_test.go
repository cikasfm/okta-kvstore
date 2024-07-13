package server

import (
	"codesignal/internal/api"
	"codesignal/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_SetupRoutes(t *testing.T) {

	keyValueStoreApi := api.NewKeyValueStoreApi(services.NewInMemoryKeyValueStore())

	type testrequest struct {
		method  string
		uri     string
		content string
	}
	tests := []struct {
		name       string
		req        testrequest
		wantStatus int
		wantBody   string
	}{
		{
			name: "Test get key on empty",
			req: testrequest{
				method:  http.MethodGet,
				uri:     "/key/testKey",
				content: "",
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "{\"message\":\"key testKey not found\"}",
		},
		{
			name: "Test set key",
			req: testrequest{
				method:  http.MethodPost,
				uri:     "/key",
				content: "{\"testKey\":\"testValue\"}",
			},
			wantStatus: http.StatusOK,
			wantBody:   "{\"testKey\":\"testValue\"}",
		},
		{
			name: "Test get key after set",
			req: testrequest{
				method:  http.MethodGet,
				uri:     "/key/testKey",
				content: "",
			},
			wantStatus: http.StatusOK,
			wantBody:   "{\"testKey\":\"testValue\"}",
		},
		{
			name: "Test delete key",
			req: testrequest{
				method:  http.MethodDelete,
				uri:     "/key/testKey",
				content: "",
			},
			wantStatus: http.StatusOK,
			wantBody:   "{\"message\":\"key deleted successfully\"}",
		},
		{
			name: "Test get key after delete",
			req: testrequest{
				method:  http.MethodGet,
				uri:     "/key/testKey",
				content: "",
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "{\"message\":\"key testKey not found\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			ctx, engine := gin.CreateTestContext(w)
			SetupRoutes(engine, keyValueStoreApi)

			req, err := http.NewRequestWithContext(ctx, tt.req.method, tt.req.uri, strings.NewReader(tt.req.content))

			assert.NoError(t, err)

			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())

		})
	}
}
