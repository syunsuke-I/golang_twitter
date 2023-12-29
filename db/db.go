package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

const tableNameUser = "users"

type Database struct {
	DB *gorm.DB
}

func NewDatabase() *Database {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Connected to DB successfully!")

	return &Database{DB: db}
}

func (d *Database) CreateTables() error {
	cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    password TEXT,
    activation_token VARCHAR(64),
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP
	)`, tableNameUser)

	result := d.DB.Exec(cmdU)
	if result.Error != nil {
		return result.Error // エラーを返す
	}

	return nil // 成功の場合は nil を返す
}

// Close はデータベース接続を閉じます。
func (d *Database) Close() {
	sqlDB, _ := d.DB.DB()
	sqlDB.Close()
}
