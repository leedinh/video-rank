package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func HandleInteraction(c *gin.Context, r *redis.Client) {
	// Get the user ID from the request
	userID := c.Param("userID")

	// Get the user's rank from Redis
	rank, err := r.Get(context.Background(), userID).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting user rank",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rank": rank,
	})
}
