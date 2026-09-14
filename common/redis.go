package common

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/modelbus/one-api-pro/common/logger"
)

var RDB redis.Cmdable
var RedisEnabled = true

// InitRedisClient This function is called after init()
func InitRedisClient() (err error) {
	if os.Getenv("REDIS_CONN_STRING") == "" {
		RedisEnabled = false
		logger.SysLog("REDIS_CONN_STRING not set, Redis is not enabled")
		return nil
	}
	if os.Getenv("SYNC_FREQUENCY") == "" {
		RedisEnabled = false
		logger.SysLog("SYNC_FREQUENCY not set, Redis is disabled")
		return nil
	}
	redisConnString := os.Getenv("REDIS_CONN_STRING")
	if os.Getenv("REDIS_MASTER_NAME") == "" {
		logger.SysLog("Redis is enabled")
		opt, err := redis.ParseURL(redisConnString)
		if err != nil {
			logger.FatalLog("failed to parse Redis connection string: " + err.Error())
		}
		RDB = redis.NewClient(opt)
	} else {
		// cluster mode
		logger.SysLog("Redis cluster mode enabled")
		RDB = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:      strings.Split(redisConnString, ","),
			Password:   os.Getenv("REDIS_PASSWORD"),
			MasterName: os.Getenv("REDIS_MASTER_NAME"),
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = RDB.Ping(ctx).Result()
	if err != nil {
		logger.FatalLog("Redis ping test failed: " + err.Error())
	}
	return err
}

func ParseRedisOption() *redis.Options {
	opt, err := redis.ParseURL(os.Getenv("REDIS_CONN_STRING"))
	if err != nil {
		logger.FatalLog("failed to parse Redis connection string: " + err.Error())
	}
	return opt
}

func RedisSet(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	return RDB.Set(ctx, key, value, expiration).Err()
}

func RedisGet(key string) (string, error) {
	ctx := context.Background()
	return RDB.Get(ctx, key).Result()
}

func RedisDel(key string) error {
	ctx := context.Background()
	return RDB.Del(ctx, key).Err()
}

func RedisDecrease(key string, value int64) error {
	ctx := context.Background()
	return RDB.DecrBy(ctx, key, value).Err()
}

// RedisIncrease adds value to the integer stored at key (atomic INCRBY).
// Used to roll back a previous RedisDecrease when the corresponding DB write
// fails, so the cache does not drift further from users.quota.
// RedisIncrease 将 value 加到 key 存储的整数上（原子 INCRBY）。
// 用于对应的 DB 写入失败时回滚此前已执行的 RedisDecrease，
// 避免缓存进一步偏离 users.quota。
func RedisIncrease(key string, value int64) error {
	ctx := context.Background()
	return RDB.IncrBy(ctx, key, value).Err()
}
