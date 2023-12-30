package utils

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// Connection
func RedisConnection() redis.Conn {
	const Addr = "redis:6379"

	c, err := redis.Dial("tcp", Addr)
	if err != nil {
		panic(err)
	}
	return c
}

type Data struct {
	Key   string
	Value string
}

// 複数のデータの登録(Redis: MSET key [key...])
func Mset(data []Data, c redis.Conn) {
	var query []interface{}
	for _, v := range data {
		query = append(query, v.Key, v.Value)
	}
	fmt.Println(query) // [key1 value1 key2 value2]

	c.Do("MSET", query...)
}

// 複数の値を取得 (Redis: MGET key [key...])
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
