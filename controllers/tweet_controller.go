package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/syunsuke-I/golang_twitter/models"
	"github.com/syunsuke-I/golang_twitter/utils"
)

func TweetCreate(c *gin.Context) {

	repo := models.NewRepository(db.DB)

	content := c.Request.FormValue("content")
	errMsg, err := models.LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	uid, err := c.Cookie("uid")
	if err != nil {
		// セッションIDが見つからない場合はエラーハンドリング
		fmt.Println(errMsg.SessionInvalid)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.SessionInvalid})
		return
	}

	userId, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		// strconv.ParseUintからのエラーがある場合、エラーレスポンスを返す
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なユーザーIDです。"})
		return
	}

	if userId == 0 {
		// ユーザーIDが0の場合、セッションが存在しないと見なし、エラーレスポンスを返す
		c.JSON(http.StatusBadRequest, gin.H{"error": "セッションが存在しません。"})
		return
	}

	tweet := models.Tweet{
		UserID:  userId,
		Content: content,
	}

	entry, errorMessages := repo.CreateTweet(&tweet)

	form, _ := c.MultipartForm()
	var files []*multipart.FileHeader
	for key, fileHeaders := range form.File {
		if key == "images[]" || strings.HasPrefix(key, "images[") && strings.HasSuffix(key, "]") {
			files = append(files, fileHeaders...)
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()

			url, err := utils.UploadImg(file, c)
			if err != nil {
				errChan <- err
				return
			}
			fmt.Printf("Uploaded File: %s, Size: %d\n", file.Filename, file.Size)

			imgUrl := models.Image{
				ImgUrl:  url,
				TweetID: entry.ID,
			}

			_, err = repo.CreateImage(&imgUrl)
			if err != nil {
				errChan <- err
				return
			}
			errChan <- nil
		}(file)
	}

	wg.Wait()
	close(errChan)

	// エラーチェック
	for err := range errChan {
		if err != nil {
			c.String(http.StatusInternalServerError, "Save uploaded file err: %s", err.Error())
			return
		}
	}

	// エラーがある場合はエラーメッセージを返す
	if errorMessages != nil {
		messages := []string{errorMessages.Error()}
		c.HTML(http.StatusBadRequest, "home/home.html", gin.H{
			"errorMessages": messages,
			"Tweet":         tweet,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ポストを送信しました",
	})
}
