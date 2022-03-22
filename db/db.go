package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/kuritaeiji/todo-gin-back/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db       *gorm.DB
	database *sql.DB
	err      error
)

func Init() {
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	logLevelInt, _ := strconv.Atoi(os.Getenv("MYSQL_LOG_LEVEL"))
	initDB(dsn, logger.LogLevel(logLevelInt))
}

func initDB(dsn string, logLevel logger.LogLevel) {
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to open mysql\n%v", err.Error()))
	}
	migrate()
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	database, err = db.DB()
	if err == nil {
		database.Close()
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to close mysql\n%v", err.Error()))
	}
}

func migrate() {
	db.AutoMigrate(model.User{})
}

// test
func DeleteAll() {
	db.Exec("DELETE FROM users")
}
