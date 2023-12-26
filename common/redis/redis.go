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

func NewRedis(host, pass string, idtime, maxid, maxac int) (*NeRedis, error) {
	REDIS = &NeRedis{}

	// 创建Redis连接池
	REDIS.Conn = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass, // 如果有密码，可以在这里设置
		DB:       0,    // 选择要使用的Redis数据库
		PoolSize: 10,   // 连接池的大小
	})

	// 连接测活
	_, err := REDIS.Conn.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	REDIS.Status = true
	fmt.Println("连接Redis成功")
	REDIS.Status = true
	return REDIS, nil
}

// 设置缓存
func (a *NeRedis) Set(key string, data interface{}, ts time.Duration) error {
	// conn := a.Conn.Get()
	// defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// _, err = conn.Do("SET", key, value)
	// if err != nil {
	// 	return err
	// }

	// _, err = conn.Do("EXPIRE", key, time)
	// if err != nil {
	// 	return err
	// }
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
