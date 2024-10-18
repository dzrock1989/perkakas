package rds

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Rediser interface {
	redisRateLimiter

	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd

	Close() error
}

type redisRateLimiter interface {
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
	ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd
	ScriptLoad(ctx context.Context, script string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd

	EvalRO(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	EvalShaRO(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
}
