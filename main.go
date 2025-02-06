package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/leedinh/video-rank/handlers"
)

var redisClient *redis.Client

func initRedis() {
	// Connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func main() {
	initRedis()
	r := gin.Default()
	r.Group("/api")

	r.POST("/interactions", func(c *gin.Context) {
		handlers.HandleInteraction(c, redisClient)
	})

	r.GET("/ranks/:userID", func(c *gin.Context) {
		handlers.HandleGetRank(c, redisClient)
	})

	r.Run(":8080")
}
