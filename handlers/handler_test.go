package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestHandlerInteractionAPI(t *testing.T) {
	ResetRedis()
	router := gin.Default()
	router.POST("/api/interactions", func(ctx *gin.Context) {
		HandleInteraction(ctx, RedisClient)
	})

	// Test case 1: Add interaction
	w := httptest.NewRecorder()
	reqBody := `{"user_id": "user1", "video_id": "video1", "interaction_type": "view", "value": 1}`
	req := httptest.NewRequest("POST", "/api/interactions", strings.NewReader(reqBody))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test case 2: Invalid interaction type
	w = httptest.NewRecorder()
	reqBody = `{"user_id": "user1", "video_id": "video1", "interaction_type": "invalid", "value": 1}`
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestGetRankings(t *testing.T) {
	ResetRedis()
	router := gin.Default()
	router.GET("/api/rankings", func(ctx *gin.Context) {
		GetRankings(ctx, RedisClient)
	})

	// Test case 1: No rankings
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/rankings", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `[]`, w.Body.String())

	// Test case 2: Add some rankings
	RedisClient.ZAdd(ctx, "global_rank", &redis.Z{
		Score:  10,
		Member: "video1",
	})
	RedisClient.ZAdd(ctx, "global_rank", &redis.Z{
		Score:  20,
		Member: "video2",
	})
	RedisClient.ZAdd(ctx, "global_rank", &redis.Z{
		Score:  30,
		Member: "video3",
	})
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/rankings", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `[{"video_id":"video3","score":30},{"video_id":"video2","score":20},{"video_id":"video1","score":10}]`, w.Body.String())
}
