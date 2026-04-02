package database

import (
	"context"
	"log"

	"email-dispatcher/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.GetEnv("REDIS_URL"),
		Password: config.GetEnv("REDIS_PASSWORD"), // password from env
		DB:       0,                               // use default DB
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Connected to Redis")
}
