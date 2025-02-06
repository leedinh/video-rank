package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/leedinh/video-rank/docs"
	"github.com/leedinh/video-rank/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Video Rank API
// @version 1.0
// @description This is a Microservice Ranking Video to rank videos based on user interactions.
// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("/api")
	{
		v1.POST("/interactions", func(c *gin.Context) {
			handlers.HandleInteraction(c, redisClient)
		})

		v1.GET("/rankings", func(c *gin.Context) {
			handlers.GetRankings(c, redisClient)
		})
	}
	r.Run(":8080")
}
