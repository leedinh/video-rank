package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// @Summary Get rankings
// @Description Get the top-n global rankings or user rankings
// @Tags rankings
// @Accept json
// @Produce json
// @Param user_id query string false "User ID for personalized ranking"
// @Param limit query int false "Number of results to return" default(10)
// @Success 200 {array} VideoScore
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /rankings [get]
func GetRankings(c *gin.Context, r *redis.Client) {
	userID := c.Query("user_id")
	limitStr := c.DefaultQuery("limit", "10")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit",
		})
		return
	}

	key := "global_rank"
	if userID != "" {
		key = "user:" + userID + ":rank"
	}

	rankings, err := r.ZRevRangeWithScores(c, key, 0, int64(limit-1)).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get rankings",
		})
		return
	}

	response := make([]VideoScore, len(rankings))
	for i, value := range rankings {
		response[i] = VideoScore{
			VideoID: value.Member.(string),
			Score:   value.Score,
		}
	}

	c.JSON(http.StatusOK, response)
}

type VideoScore struct {
	VideoID string  `json:"video_id"`
	Score   float64 `json:"score"`
}
