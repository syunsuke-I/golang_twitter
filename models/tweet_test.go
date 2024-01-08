package models

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ツイート投稿　正常系
func TestCreateTweet(t *testing.T) {
	os.Chdir("..")                 // プロジェクトのルートに移動する
	db, mock, err := sqlmock.New() // モックデータベース接続
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	uid := "1"
	userId, _ := strconv.ParseUint(uid, 10, 64)
	content := "test message"

	// Mock設定
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tweets" ("user_id","content","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(userId, content, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// 実行
	repo := NewRepository(gdb)

	user, err := repo.CreateTweet(&Tweet{UserID: userId, Content: content})
	fmt.Println("user = ", user)
	fmt.Println("err = ", err)

	// モックが期待通りの動作をしたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestCreateTweetMaxLength(t *testing.T) {
	os.Chdir("..")                 // プロジェクトのルートに移動する
	db, mock, err := sqlmock.New() // モックデータベース接続
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	uid := "1"
	userId, _ := strconv.ParseUint(uid, 10, 64)
	content := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	if len(content) == 140 {
		fmt.Println(len(content))
	} else {
		fmt.Println(len(content))
		content = ""
	}

	// Mock設定
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tweets" ("user_id","content","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(userId, content, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// 実行
	repo := NewRepository(gdb)

	user, err := repo.CreateTweet(&Tweet{UserID: userId, Content: content})
	fmt.Println("user = ", user)
	fmt.Println("err = ", err)

	// モックが期待通りの動作をしたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

// ツイート投稿　異常系
func TestCreateTweetRequired(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	uid := "1"
	userId, _ := strconv.ParseUint(uid, 10, 64)
	content := ""

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateTweet(&Tweet{UserID: userId, Content: content})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "Content: " + errMsg.TweetRequired + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateTweetTooLong(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	uid := "1"
	userId, _ := strconv.ParseUint(uid, 10, 64)
	content := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	if len(content) == 141 {
		fmt.Println(len(content))
	} else {
		fmt.Println(len(content))
		content = ""
	}

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateTweet(&Tweet{UserID: userId, Content: content})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "Content: " + errMsg.TweetLength + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}
