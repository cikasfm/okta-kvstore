package main

import (
	"codesignal/cmd/store/services"
	"encoding/json"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type ErrorMessage struct {
	message string
}

type HttpHandler struct {
	storage services.IKeyValueStore
}

func NewHttpHandler(service services.IKeyValueStore) *HttpHandler {
	return &HttpHandler{
		storage: service,
	}
}

// GET /key/:key
// If the key exists, return the key-value pair, for example {"key":"value"}.
// If the key does not exist, return a not found message, for example {"message":"key not found"}.
func (h *HttpHandler) getByKey(c *gin.Context) {
	key := c.Query("key")
	if key != "" {
		value, err := h.storage.Get(key)
		if err != nil {
			status, message := handleError(err)
			c.JSON(status, message)
		} else if value != "" {
			c.JSON(http.StatusOK, gin.H{key: value})
		} else {
			c.JSON(http.StatusNotFound, ErrorMessage{spew.Sprintf("key '%s' not found", key)})
		}
	} else {
		c.JSON(http.StatusBadRequest, ErrorMessage{"parameter 'key' is required"})
	}
}

// POST /key
// Should accept a payload like {"key":"value"} and return a message indicating if the operation was successful.
// If the key does not exist, create the key-value pair and return a success message, for example {"message":"key created successfully"}.
// If the key already exists, return a conflict message, for example {"message":"key already exists"}.
func (h *HttpHandler) setValue(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		status, message := handleError(err)
		c.JSON(status, message)
	}
	var data = map[string]string{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		status, message := handleError(err)
		c.JSON(status, message)
		return
	}
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, ErrorMessage{"no data provided"})
		return
	}
	if len(data) > 1 {
		c.JSON(http.StatusBadRequest, ErrorMessage{"only a single key can be accepted"})
		return
	}
	for key, value := range data {
		err = h.storage.Set(key, value)
		if err != nil {
			status, message := handleError(err)
			c.JSON(status, message)
		} else {
			c.JSON(http.StatusOK, gin.H{})
		}
	}
}

// DELETE /key/:key
// If the key exists, delete the key and return a success message, for example {"message":"key deleted successfully"}.
// If the key does not exist, return a not found message, for example {"message":"key not found"}.
func (h *HttpHandler) deleteValue(c *gin.Context) {
	key := c.Query("key")
	if key != "" {
		err := h.storage.Delete(key)
		if err != nil {
			status, message := handleError(err)
			c.JSON(status, message)
		} else {
			c.JSON(http.StatusOK, ErrorMessage{"key deleted successfully"})
		}
	} else {
		c.JSON(http.StatusBadRequest, ErrorMessage{"parameter 'key' is required"})
	}
}

func handleError(err error) (int, string) {
	var serviceError services.ServiceError
	if errors.As(err, &serviceError) {
		switch serviceError.Code {
		case services.CodeKeyNotFound:
			return http.StatusNotFound, serviceError.Message
		case services.CodeKeyExists:
			return http.StatusConflict, serviceError.Message
		default:
			return http.StatusInternalServerError, serviceError.Message
		}
	}
	return http.StatusInternalServerError, err.Error()
}
