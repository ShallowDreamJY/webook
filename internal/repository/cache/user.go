package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/internal/domain"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	// 传单机 Redis 或者 cluster Redis 都可以
	client     redis.Cmdable
	expiration time.Duration
}

// A 用到了 B ，B 一定是接口 -》保证面向接口
// A 用到了 B ，B 一定是 A 的字段 -》规避包变量，包方法
// A 用到了 B ，A 绝对不初始化 B ，而是外面注入
func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 60,
	}
}

// 只要 error 为 nil ，就认为缓存里有数据
// 如果没有数据则返回特定的 error
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, nil
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)

	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
