package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/syunsuke-I/golang_twitter/models"
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
func SessionAuthMiddleware(redisClient redis.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		errMsg, err := models.LoadConfig("settings/error_messages.json")
		if err != nil {
			fmt.Println("Error loading config:", err)
		}

		uid, err := c.Cookie("uid")
		if err != nil {
			// セッションIDが見つからない場合はエラーハンドリング
			c.HTML(http.StatusBadRequest, "login/login.html", gin.H{
				"errorMessages": []string{errMsg.LoginRequired},
			})
			return
		}

		// ここでredisとの照会を行う
		sid, err := GetSessionIDByUserID(uid, redisClient)
		if err != nil {
			c.HTML(http.StatusBadRequest, "login/login.html", gin.H{
				"errorMessages": []string{errMsg.SessionInvalid},
			})
			return
		}

		err = IsSessionValid(sid, redisClient)
		if err != nil {
			c.HTML(http.StatusBadRequest, "login/login.html", gin.H{
				"errorMessages": []string{errMsg.SessionInvalid},
			})
			return
		}
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

func GetSessionIDByUserID(userID string, c redis.Conn) (string, error) {
	userSessionKey := "user:" + userID + ":session"

	// RedisからセッションIDを取得
	sessionID, err := redis.String(c.Do("GET", userSessionKey))
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func setKeyWithExpiration(c redis.Conn, key string, value string, ttl time.Duration) error {
	_, err := c.Do("SET", key, value, "EX", int(ttl.Seconds()))
	if err != nil {
		return err
	}
	return nil
}

func IsSessionValid(sessionID string, c redis.Conn) error {
	sessionKey := "session:" + sessionID

	ttl, err := redis.Int(c.Do("TTL", sessionKey))
	if err != nil {
		return err
	}

	// 有効であればセッションの期限を更新する
	if ttl > 0 {
		refreshSession(sessionID, c)
		return nil
	}

	// それ以外の場合はセッションは無効
	return errors.New("session expired or does not exist")
}

// セッションの有効期限を更新する
func refreshSession(sessionID string, c redis.Conn) error {
	_, err := c.Do("EXPIRE", "session:"+sessionID, 1800) // 1800秒（30分）に設定
	return err
}
