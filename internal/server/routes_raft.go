package server

import (
	"codesignal/internal/api"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/raft"
	"log"
	"net/http"
)

func JoinHandler(ra *raft.Raft) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NodeID  string `json:"node_id"`
			Address string `json:"address"`
		}

		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			c.JSON(http.StatusBadRequest, api.ErrorMessage{Message: "Invalid request payload"})
			return
		}

		configFuture := ra.GetConfiguration()
		if err := configFuture.Error(); err != nil {
			c.JSON(http.StatusInternalServerError, api.ErrorMessage{Message: "Failed to get raft configuration"})
			return
		}

		for _, srv := range configFuture.Configuration().Servers {
			if srv.ID == raft.ServerID(req.NodeID) || srv.Address == raft.ServerAddress(req.Address) {
				c.JSON(http.StatusBadRequest, api.ErrorMessage{Message: "Node already exists"})
				return
			}
		}

		future := ra.AddVoter(raft.ServerID(req.NodeID), raft.ServerAddress(req.Address), 0, 0)
		if err := future.Error(); err != nil {
			log.Println(fmt.Errorf("AddVoter error: %v", err))
			c.JSON(http.StatusInternalServerError, api.ErrorMessage{Message: "Failed to add node to cluster"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Added node to cluster"})
	}
}
