package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type NeRedis struct {
	Conn   *redis.Client
	Status bool
}

var REDIS *NeRedis
var ctx = context.Background()

// NewRedis 创建 Redis 客户端，参数支持默认值和校验
func NewRedis(host, pass string, idleTimeout, maxIdle, maxActive, minIdle int) (*NeRedis, error) {
	if host == "" {
		return nil, fmt.Errorf("Redis地址不能为空")
	}
	if idleTimeout <= 0 {
		idleTimeout = 300
	}
	if maxIdle <= 0 {
		maxIdle = 10
	}
	if maxActive <= 0 {
		maxActive = 100
	}
	if minIdle < 0 {
		minIdle = 0
	}

	opts := &redis.Options{
		Addr:            host,
		Password:        pass,
		DB:              0,
		DialTimeout:     10 * time.Second,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		ConnMaxIdleTime: time.Duration(idleTimeout) * time.Second, // go-redis v9 推荐
		MaxIdleConns:    maxIdle,
		PoolSize:        maxActive,
		PoolTimeout:     5 * time.Second,
		MinIdleConns:    minIdle,
	}

	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		fmt.Printf("Redis连接失败: %v, Addr: %s\n", err, host)
		return nil, fmt.Errorf("Redis连接失败: %v", err)
	}

	REDIS = &NeRedis{Conn: client, Status: true}
	fmt.Printf("Redis连接成功: %s\n", host)
	return REDIS, nil
}

// Close 优雅关闭 Redis 连接
func (a *NeRedis) Close() error {
	if a.Conn != nil {
		return a.Conn.Close()
	}
	return nil
}

// Set 设置缓存，ts为秒
func (a *NeRedis) Set(key string, data interface{}, ts time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return a.Conn.Set(ctx, key, value, ts*time.Second).Err()
}

// Exists 检查缓存中是否存在 key
func (a *NeRedis) Exists(key string) bool {
	res, err := a.Conn.Exists(ctx, key).Result()
	if err != nil {
		fmt.Printf("Redis Exists失败: %v, key: %s\n", err, key)
		return false
	}
	return res > 0
}

// Get 获取缓存数据并反序列化到 v
func (a *NeRedis) Get(key string, v interface{}) error {
	reply, err := a.Conn.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(reply, v)
}

// Delete 删除缓存 key
func (a *NeRedis) Delete(key string) (int64, error) {
	return a.Conn.Del(ctx, key).Result()
}

// SetLong 永久保存 key，无过期时间
func (a *NeRedis) SetLong(key string, data interface{}) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return a.Conn.Set(ctx, key, value, 0).Err()
}
