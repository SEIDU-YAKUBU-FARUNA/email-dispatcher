package database

import (
	"time"
)

// Limit: 5 requests per minute per user
func AllowRequest(key string) (bool, error) {

	count, err := RedisClient.Incr(Ctx, key).Result()
	if err != nil {
		return false, err
	}

	// First request → set expiration
	if count == 1 {
		RedisClient.Expire(Ctx, key, time.Minute)
	}

	if count > 5 {
		return false, nil
	}

	return true, nil
}
