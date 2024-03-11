package redis

import (
	"context"
	"errors"
	"fmt"
	"go-test/internal/config"
	"go-test/internal/log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var (
	client *redis.Client
	rs     *redsync.Redsync
)

// ConnectToRedis 连接到Redis服务器
func ConnectToRedis(c *config.RedisConfig) error {
	client = redis.NewClient(&redis.Options{
		Addr:        c.Addr,
		Password:    c.Password,
		DB:          c.Db,
		IdleTimeout: 240 * time.Second,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	rs = redsync.New(goredis.NewPool(client))

	log.Infof("connected to redis, addr: %s", c.Addr)
	return nil
}

// Get 操作命令
func Get(key string) (string, error) {
	return client.Get(context.Background(), key).Result()
}

// Set 操作命令
func Set(key string, value interface{}) error {
	return client.Set(context.Background(), key, value, 0).Err()
}

// Del 操作命令
func Del(key string) error {
	return client.Del(context.Background(), key).Err()
}

// Incr 操作命令
func Incr(key string) (int64, error) {
	return client.Incr(context.Background(), key).Result()
}

// Exists 操作命令
func Exists(key string) bool {
	return client.Exists(context.Background(), key).Val() > 0
}

// Expire 操作命令
func Expire(key string, expire time.Duration) (bool, error) {
	return client.Expire(context.Background(), key, expire).Result()
}

// SetEx 操作命令
func SetEX(key string, value interface{}, expire time.Duration) error {
	return client.SetEX(context.Background(), key, value, expire).Err()
}

// SetNx 操作命令
func SetNX(key string, value interface{}, expire time.Duration) (bool, error) {
	return client.SetNX(context.Background(), key, value, expire).Result()
}

// HGet 操作命令
func HGet(key string, field string) (string, error) {
	return client.HGet(context.Background(), key, field).Result()
}

// HSet 操作命令
func HSet(key string, field string, value interface{}) error {
	return client.HSet(context.Background(), key, field, value).Err()
}

// HGetAll 操作命令
func HGetAll(key string) (map[string]string, error) {
	return client.HGetAll(context.Background(), key).Result()
}

// HKeys 操作命令
func HKeys(key string) ([]string, error) {
	return client.HKeys(context.Background(), key).Result()
}

// HIncrBy 操作命令
func HIncrBy(key string, field string, value int64) (int64, error) {
	return client.HIncrBy(context.Background(), key, field, value).Result()
}

// LRange 获取指定范围内的数据列表
func LRange(key string, start, end int) (interface{}, error) {
	return client.LRange(context.Background(), key, int64(start), int64(end)).Result()
}

// LPush 向列表的前面添加数据
func LPush(key string, value ...interface{}) error {
	return client.LPush(context.Background(), key, value).Err()
}

// RPop 移除列表最后面的数据
func RPop(key string) (string, error) {
	return client.RPop(context.Background(), key).Result()
}

// LLen 返回列表的长度
func LLen(key string) (int64, error) {
	return client.LLen(context.Background(), key).Result()
}

// ZAdd 向列表添加数据
func ZAdd(key string, value string, score float64) (int64, error) {
	return client.ZAdd(context.Background(), key, &redis.Z{Score: score, Member: value}).Result()
}

// ZRem 移除指定的数据
func ZRem(key string, value string) error {
	return client.ZRem(context.Background(), key, value).Err()
}

// ZCard 返回列表的长度
func ZCard(key string) (int64, error) {
	return client.ZCard(context.Background(), key).Result()
}

// ZIncrBy 向列表的指定字段增加分数(不存在时会创建)
func ZIncrBy(key string, value string, score float64) error {
	return client.ZIncrBy(context.Background(), key, score, value).Err()
}

// ZRange 获取指定索引范围内的成员数据列表(正序)
func ZRange(key string, start, end int) ([]string, error) {
	return client.ZRange(context.Background(), key, int64(start), int64(end)).Result()
}

// ZRevRange 获取指定索引范围内的成员数据列表(倒序)
func ZRevRange(key string, start, end int) ([]string, error) {
	return client.ZRevRange(context.Background(), key, int64(start), int64(end)).Result()
}

// ZRangeWithScores 获取指定索引范围内的成员与分数列表(正序)
func ZRangeWithScores(key string, start, end int) ([]redis.Z, error) {
	return client.ZRangeWithScores(context.Background(), key, int64(start), int64(end)).Result()
}

// ZRevRangeWithScores 获取指定索引范围内的成员与分数列表(倒序)
func ZRevRangeWithScores(key string, start, end int) ([]redis.Z, error) {
	return client.ZRevRangeWithScores(context.Background(), key, int64(start), int64(end)).Result()
}

// ZRangeByScore 获取指定分数范围内的成员数据列表
func ZRangeByScore(key string, start, end int) ([]string, error) {
	return client.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{Min: strconv.Itoa(start), Max: strconv.Itoa(end)}).Result()
}

// ZRank 获取指定成员的排名
func ZRank(key, idStr string) (int64, error) {
	result, err := client.ZRank(context.Background(), key, idStr).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return result, nil
}

// ZRevRank 获取指定成员的排名(倒序)
func ZRevRank(key, idStr string) (int64, error) {
	result, err := client.ZRevRank(context.Background(), key, idStr).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return result, nil
}

// ZScore 获取指定成员的分数
func ZScore(key, idStr string) (string, error) {
	result, err := client.ZScore(context.Background(), key, idStr).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return strconv.FormatFloat(result, 'f', 2, 64), err
}

// NewLock 获取分布式锁
func NewLock(key string) *redsync.Mutex {
	return rs.NewMutex(key)
}

// NewDistributedLock 获取分布式锁
func NewDistributedLock() *DistributedLock {
	return &DistributedLock{
		rs: rs,
	}
}

func (dl *DistributedLock) Lock(resource string) error {
	mutex := dl.rs.NewMutex(resource, redsync.WithExpiry(20*time.Second))
	if err := mutex.Lock(); err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	dl.mutex = mutex
	dl.locked = true
	return nil
}

func (dl *DistributedLock) Unlock() error {
	if !dl.locked {
		return nil
	}

	if ok, err := dl.mutex.Unlock(); err != nil || !ok {
		return fmt.Errorf("failed to release lock: %w", err)
	}
	dl.locked = false
	return nil
}

type DistributedLock struct {
	rs     *redsync.Redsync
	mutex  *redsync.Mutex
	locked bool
}
