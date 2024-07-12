package api

import (
	"codesignal/internal/services"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

// IKeyValueStoreApi handles Get / Set / Delete HTTP REST API calls
type IKeyValueStoreApi interface {
	GetByKey(c *gin.Context)
	SetValue(c *gin.Context)
	DeleteValue(c *gin.Context)
}

type KeyValueStoreApi struct {
	storage services.IKeyValueStore
}

func NewKeyValueStoreApi(service services.IKeyValueStore) IKeyValueStoreApi {
	return &KeyValueStoreApi{
		storage: service,
	}
}

// GetByKey GET /key/:key
//
// If the key exists, return the key-value pair, for example {"key":"value"}.
//
// If the key does not exist, return a not found message, for example: {"message":"key not found"}.
func (h *KeyValueStoreApi) GetByKey(c *gin.Context) {
	key := c.Param("key")
	if key != "" {
		value, err := h.storage.Get(key)
		if err != nil {
			status, message := handleError(err)
			c.JSON(status, ErrorMessage{message})
		} else if value != "" {
			c.JSON(http.StatusOK, gin.H{key: value})
		} else {
			c.JSON(http.StatusNotFound, ErrorMessage{fmt.Sprintf("key '%s' not found", key)})
		}
	} else {
		c.JSON(http.StatusBadRequest, ErrorMessage{"parameter 'key' is required"})
	}
}

// SetValue POST /key
//
// Should accept a payload like {"key":"value"} and return a message indicating if the operation was successful.
//
// If the key does not exist, create the key-value pair and return a success message, for example {"message":"key created successfully"}.
//
// If the key already exists, return a conflict message, for example {"message":"key already exists"}.
func (h *KeyValueStoreApi) SetValue(c *gin.Context) {
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{"empty request body"})
		return
	}
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{"reading request body failed"})
		return
	}
	var data = map[string]string{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{"request is not a valid json"})
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
			c.JSON(status, ErrorMessage{message})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				key: value,
			})
			return
		}
	}
}

// DeleteValue DELETE /key/:key
//
// If the key exists, delete the key and return a success message, for example {"message":"key deleted successfully"}.
//
// If the key does not exist, return a not found message, for example {"message":"key not found"}.
func (h *KeyValueStoreApi) DeleteValue(c *gin.Context) {
	key := c.Param("key")
	if key != "" {
		err := h.storage.Delete(key)
		if err != nil {
			status, message := handleError(err)
			c.JSON(status, ErrorMessage{message})
		} else {
			c.JSON(http.StatusOK, ErrorMessage{"key deleted successfully"})
		}
	} else {
		// seems to be an impossible case
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
