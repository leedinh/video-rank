package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type InteractionRequest struct {
	UserID          string  `json:"user_id"`
	VideoID         string  `json:"video_id" binding:"required"`
	InteractionType string  `json:"interaction_type" binding:"required"`
	Value           float64 `json:"value"`
}

// @Summary Handle interaction
// @Description Update the rank of a video based on user interaction
// @Tags interactions
// @Accept json
// @Produce json
// @Param interaction body InteractionRequest true "Interaction details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /interactions [post]
func HandleInteraction(c *gin.Context, r *redis.Client) {
	var interaction InteractionRequest
	if err := c.ShouldBindBodyWithJSON(&interaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// Calculate score based on interaction type
	score, err := calculateScore(interaction.InteractionType, interaction.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid interaction type",
		})
		return
	}

	// If user ID is not provided, update global rank
	if interaction.UserID != "" {
		if err := r.ZIncrBy(c, "user:"+interaction.UserID+":rank", score, interaction.VideoID).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update rank",
			})
			return
		}
	}

	if err := r.ZIncrBy(c, "global_rank", score, interaction.VideoID).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update global rank",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
	})
}

// Assume score formula is as follows:
func calculateScore(interactionType string, value float64) (float64, error) {
	switch interactionType {
	case "view":
		return 1, nil
	case "like":
		return 2, nil
	case "comment":
		return 3, nil
	case "share":
		return 4, nil
	case "watch_time":
		return value * 0.5, nil
	default:
		return 0, fmt.Errorf("invalid interaction type")
	}
}
