package models

import (
	"database/sql"
	"fmt"
	"log"
)

var Db *sql.DB
var err error

const (
	host          = "postgres"
	port          = 5432
	user          = "postgres"
	password      = "password"
	dbname        = "postgres"
	tableNameUser = "users"
)

func Init() {
	// 接続文字列を作成
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// PostgreSQLに接続
	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// 接続を確認
	err = Db.Ping()
	if err != nil {
		log.Fatal("接続失敗:", err)
	}

	fmt.Println("接続成功!")

}
