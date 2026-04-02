package ratelimiter

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redis     *redis.Client
	maxReqs   int
	windowSec int
}

func NewRateLimiter(redisClient *redis.Client, maxRequests int, windowSeconds int) *RateLimiter {
	return &RateLimiter{
		redis:     redisClient,
		maxReqs:   maxRequests,
		windowSec: windowSeconds,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, int, error) {
	now := time.Now().UnixMilli()
	windowStart := now - int64(rl.windowSec)*1000

	pipe := rl.redis.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "-inf", strconv.FormatInt(windowStart, 10))
	cardCmd := pipe.ZCard(ctx, key)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, err
	}

	count := cardCmd.Val()
	remaining := rl.maxReqs - int(count)

	if int(count) >= rl.maxReqs {
		return false, 0, nil
	}

	member := fmt.Sprintf("%d:%d", now, rand.Int63())
	pipe2 := rl.redis.Pipeline()
	pipe2.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: member})
	pipe2.Expire(ctx, key, time.Duration(rl.windowSec+1)*time.Second)
	_, err = pipe2.Exec(ctx)
	if err != nil {
		return false, 0, err
	}

	return true, remaining - 1, nil
}

func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "rate_limit:" + c.ClientIP()
		allowed, remaining, err := rl.Allow(c.Request.Context(), key)
		if err != nil {
			c.JSON(500, gin.H{"error": "rate limiter error"})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.maxReqs))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if !allowed {
			c.JSON(429, gin.H{
				"error":               "too many requests",
				"retry_after_seconds": rl.windowSec,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
