package tokenblacklist

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBlacklist struct {
	redis *redis.Client
}

func NewTokenBlacklist(redisClient *redis.Client) *TokenBlacklist {
	return &TokenBlacklist{redis: redisClient}
}

func hashToken(tokenString string) string {
	h := sha256.Sum256([]byte(tokenString))
	return hex.EncodeToString(h[:])
}

func (tb *TokenBlacklist) Add(ctx context.Context, tokenString string, expiration time.Duration) error {
	key := "blacklist:" + hashToken(tokenString)
	return tb.redis.Set(ctx, key, "1", expiration).Err()
}

func (tb *TokenBlacklist) IsBlacklisted(ctx context.Context, tokenString string) (bool, error) {
	key := "blacklist:" + hashToken(tokenString)
	result, err := tb.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
