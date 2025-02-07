package handlers

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	DB:   1, // new database for testing
})

func ResetRedis() {
	RedisClient.FlushDB(ctx)
}

func TestUpdateRank(t *testing.T) {
	ResetRedis()
	// Add some initial ranks
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

	// Update score
	err := UpdateScore(RedisClient, &gin.Context{}, "global_rank", "video1", 5)
	if err != nil {
		t.Errorf("Failed to update score: %v", err)
	}

	// Check if rank is updated
	rank, err := RedisClient.ZScore(ctx, "global_rank", "video1").Result()
	if err != nil {
		t.Errorf("Failed to get rank: %v", err)
	}
	if rank != 15 {
		t.Errorf("Expected rank to be 15, got %v", rank)
	}
}
func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name            string
		interactionType string
		value           float64
		want            float64
		wantErr         bool
	}{
		{
			name:            "view",
			interactionType: "view",
			value:           1,
			want:            1,
			wantErr:         false,
		},
		{
			name:            "like",
			interactionType: "like",
			value:           1,
			want:            2,
			wantErr:         false,
		},
		{
			name:            "invalid",
			interactionType: "invalid",
			value:           1,
			want:            0,
			wantErr:         true,
		},
		{
			name:            "watch_time",
			interactionType: "watch_time",
			value:           100,
			want:            50,
			wantErr:         false,
		},

		{
			name:            "comment",
			interactionType: "comment",
			value:           1,
			want:            3,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		score, err := calculateScore(tt.interactionType, tt.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: calculateScore() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if score != tt.want {
			t.Errorf("%s: calculateScore() = %v, want %v", tt.name, score, tt.want)
		}

	}
}
