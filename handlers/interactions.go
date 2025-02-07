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
		key := fmt.Sprintf("user:%s:rank", interaction.UserID)
		if err := UpdateScore(r, c, key, interaction.VideoID, score); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update user rank: " + err.Error(),
			})
			return
		}
	}

	// Update global rank
	if err := UpdateScore(r, c, "global_rank", interaction.VideoID, score); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update global rank " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
	})
}

func UpdateScore(r *redis.Client, ctx *gin.Context, key, videoID string, score float64) error {
	if err := r.ZIncrBy(ctx, key, score, videoID).Err(); err != nil {
		return err
	}

	return nil
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
