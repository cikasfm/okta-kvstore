package server

import (
	"codesignal/internal/api"
	"codesignal/internal/store"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/raft"
	"log"
	"net/http"
)

func JoinHandler(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NodeID  string `json:"node_id"`
			Address string `json:"address"`
		}

		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			c.JSON(http.StatusBadRequest, api.ErrorMessage{Message: "Invalid request payload"})
			return
		}

		err := store.Join(req.NodeID, req.Address)
		if err != nil {
			log.Println(fmt.Sprintf("join node error: %v", err))
			if errors.Is(err, raft.ErrNotLeader) {
				c.JSON(http.StatusBadRequest,
					api.ErrorMessage{Message: "Unable to join cluster: Node is not the leader"})
			} else if errors.Is(err, raft.ErrNotVoter) {
				c.JSON(http.StatusBadRequest,
					api.ErrorMessage{Message: "Unable to join cluster: Node is not a voter"})
			} else {
				c.JSON(http.StatusInternalServerError,
					api.ErrorMessage{Message: "Unable to join cluster: Internal Server Error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Added node to cluster"})
	}
}
