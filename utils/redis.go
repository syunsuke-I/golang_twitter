package utils

import (
	"fmt"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type Data struct {
	Key   string
	Value string
}

// Connection
func RedisConnection() redis.Conn {
	const Addr = "redis:6379"

	c, err := redis.Dial("tcp", Addr)
	if err != nil {
		panic(err)
	}
	return c
}

func SetSessionWithUserID(userID uint64, c redis.Conn) error {
	// セッションIDとしてUUIDを生成
	sessionID := uuid.New().String()
	userIDStr := strconv.FormatUint(userID, 10)

	// ユーザーIDとセッションIDのペアを保存
	data := []Data{
		{Key: "session:" + sessionID, Value: userIDStr},
	}

	// Redisにセッション情報を保存
	var query []interface{}
	for _, v := range data {
		query = append(query, v.Key, v.Value)
	}
	_, err := c.Do("MSET", query...)
	if err != nil {
		return err
	}
	return nil
}

func Mget(keys []string, c redis.Conn) []string {
	var query []interface{}
	for _, v := range keys {
		query = append(query, v)
	}
	fmt.Println("MGET query:", query) // [key1 key2]

	res, err := redis.Strings(c.Do("MGET", query...))
	if err != nil {
		panic(err)
	}
	return res
}

// TTLの設定(Redis: EXPIRE key ttl)
func Expire(key string, ttl int, c redis.Conn) {
	c.Do("EXPIRE", key, ttl)
}
