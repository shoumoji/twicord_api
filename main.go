package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

// Db is database
var db *sqlx.DB

func init() {
	// TwitterAPIのURL作成と、ヘッダーへのBearer tokenの追加
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file,", err)
	}

	// MYSQLへの接続
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPass := os.Getenv("MYSQL_PASSWORD")
	mysqlProtocol := os.Getenv("MYSQL_PROTOCOL")
	mysqlAddr := os.Getenv("MYSQL_ADDRESS")
	mysqlDBName := "twicord"
	db, err = sqlx.Connect("mysql", mysqlUser+":"+mysqlPass+"@"+mysqlProtocol+"("+mysqlAddr+")"+"/"+mysqlDBName)
	if err != nil {
		log.Fatal("Error: connect MySQL,", err)
	}

	// MYSQLのスキーマ

}

func main() {
	e := echo.New()

	e.POST("/regist/twitter/:screen_name", HandleRegistByTwitterName)
	// e.GET("/registered", HandleRegistered)
	e.Start(":8000")
}
