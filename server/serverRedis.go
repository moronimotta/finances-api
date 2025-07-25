package server

import (
	"context"
	"os"

	"finances-api/db"
	logs "finances-api/utils/logs"

	"github.com/redis/go-redis/v9"
)

type RedisServer struct {
	db            db.Database
	connectionUrl string
	queueName     string
	topicName     string
	Client        *redis.Client // Export the client
}

func NewRedisServer(db db.Database) *RedisServer {
	logs.InitLogging()
	// TODO: Change Redis PORT
	connectionUrl := os.Getenv("REDIS_URL")
	if connectionUrl == "" {
		connectionUrl = "redis://localhost:6379" // Default fallback
	}

	return &RedisServer{
		db:            db,
		connectionUrl: connectionUrl,
	}
}

func (s *RedisServer) Start() *redis.Client {
	opt, err := redis.ParseURL(s.connectionUrl)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	// Test the connection
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	s.Client = rdb
	return rdb
}

// Close gracefully closes the Redis connection
func (s *RedisServer) Close() error {
	if s.Client != nil {
		return s.Client.Close()
	}
	return nil
}
