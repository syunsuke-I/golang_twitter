package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type Data struct {
	Key   string
	Value string
}

func RedisConnection() redis.Conn {
	const Addr = "redis:6379"

	c, err := redis.Dial("tcp", Addr)
	if err != nil {
		panic(err)
	}
	return c
}

// セッションが維持できているか検証する関数のスタブ
func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("hoge")
	}
}

func SetSessionWithUserID(userID uint64, c redis.Conn) error {
	// セッションIDを生成
	sessionID := uuid.New().String()
	userIDStr := strconv.FormatUint(userID, 10)

	sessionKey := "session:" + sessionID
	userSessionKey := "user:" + userIDStr + ":session"

	// Redisにセッション情報を保存
	err := setKeyWithExpiration(c, sessionKey, userIDStr, 30*time.Minute)
	if err != nil {
		return err
	}

	// ユーザーIDとセッションIDの関連付けを保存
	_, err = c.Do("SET", userSessionKey, sessionID, "EX", int(30*time.Minute.Seconds()))
	if err != nil {
		return err
	}

	return nil
}

func GetSessionIDByUserID(userID uint64, c redis.Conn) (string, error) {
	userIDStr := strconv.FormatUint(userID, 10)
	userSessionKey := "user:" + userIDStr + ":session"

	// RedisからセッションIDを取得
	sessionID, err := redis.String(c.Do("GET", userSessionKey))
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func Mget(keys []string, c redis.Conn) []string {
	var query []interface{}
	for _, v := range keys {
		query = append(query, v)
	}

	res, err := redis.Strings(c.Do("MGET", query...))
	if err != nil {
		panic(err)
	}
	return res
}

func setKeyWithExpiration(c redis.Conn, key string, value string, ttl time.Duration) error {
	_, err := c.Do("SET", key, value, "EX", int(ttl.Seconds()))
	if err != nil {
		return err
	}
	return nil
}

// セッションの有効期限を更新する
func RefreshSession(sessionID string, c redis.Conn) error {
	_, err := c.Do("EXPIRE", "session:"+sessionID, 1800) // 1800秒（30分）に設定
	return err
}

func IsSessionValid(sessionID string, c redis.Conn) error {
	sessionKey := "session:" + sessionID

	// RedisのTTLコマンドを使用して残りの存続時間を取得
	ttl, err := redis.Int(c.Do("TTL", sessionKey))
	if err != nil {
		return err
	}

	// TTLが0より大きければセッションは有効
	if ttl > 0 {
		return nil
	}

	// それ以外の場合はセッションは無効
	return errors.New("session expired or does not exist")
}
