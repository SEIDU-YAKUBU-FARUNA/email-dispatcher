package database

import (
	"time"
)

// Returns true if request already exists
func IsDuplicate(key string) (bool, error) {

	val, err := RedisClient.Get(Ctx, key).Result()

	if err == nil && val == "1" {
		return true, nil
	}

	// Store key with expiration (5 minutes)
	err = RedisClient.Set(Ctx, key, "1", 5*time.Minute).Err()
	if err != nil {
		return false, err
	}

	return false, nil
}
