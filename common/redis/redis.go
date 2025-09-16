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

func NewRedis(host, pass string, idleTimeout, maxIdle, maxActive int) (*NeRedis, error) {
	opts := &redis.Options{
		Addr:         host,
		Password:     pass,
		DB:           0,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		MaxIdleConns: maxIdle,
		PoolSize:     maxActive,
		PoolTimeout:  5 * time.Second,
		MinIdleConns: 5,
	}

	REDIS = &NeRedis{}
	REDIS.Conn = redis.NewClient(opts)

	// 使用带超时的上下文进行连接测试
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := REDIS.Conn.Ping(ctx).Result()
	if err != nil {
		worklog.Logio.WERR(fmt.Sprintf("Redis连接失败: %v, Addr: %s", err, host))
		return nil, fmt.Errorf("Redis连接失败: %v", err)
	}

	REDIS.Status = true
	worklog.Logio.WTRACE(fmt.Sprintf("Redis连接成功: %s", host))
	return REDIS, nil
}

// 设置缓存
func (a *NeRedis) Set(key string, data interface{}, ts time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	tn := ts * time.Second
	err = a.Conn.Set(ctx, key, value, tn).Err()
	if err != nil {
		fmt.Println("Failed to set key:", err)
		return err
	}
	return nil
}

// 检测缓存中是否有KEY
func (a *NeRedis) Exists(key string) bool {

	res, err := a.Conn.Do(ctx, "EXISTS", key).Bool()
	if err != nil {
		fmt.Println("rediserr", err)
		return false
	}
	return res
}

// 获取数据
func (a *NeRedis) Get(key string, v interface{}) error {

	reply, err := a.Conn.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	err = json.Unmarshal(reply, &v)
	return err
}

// 删除
func (a *NeRedis) Delete(key string) (int64, error) {

	return a.Conn.Del(ctx, key).Result()
}

func (a *NeRedis) SetLong(key string, data interface{}) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = a.Conn.Do(ctx, "SET", key, value).Err()
	if err != nil {
		return err
	}

	return nil
}
