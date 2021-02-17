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

	mysqlDSN := mysqlUser + ":" + mysqlPass + "@" + mysqlProtocol + "(" + mysqlAddr + ")" + "/" + mysqlDBName
	db, err = sqlx.Connect("mysql", mysqlDSN)
	if err != nil {
		log.Fatal("Error: connect MySQL,", err)
	}

	// MYSQLのスキーマ定義
	schema := `
	CREATE TABLE IF NOT EXISTS twitter_user(
		id BIGINT UNSIGNED NOT NULL PRIMARY KEY,
		screen_name CHAR(50) NOT NULL,
		image_url VARCHAR(2500)
	);`
	db.MustExec(schema)
}

func main() {
	e := echo.New()

	e.POST("/regist/twitter/:screen_name", HandleRegistByTwitterName)
	// フロントエンドで表にする用の今登録してるユーザを全てとってくるやつ
	// e.GET("/twitter/all", HandleRegistered)

	// memo: twitterAPIv2を使えばstreamとしてTweetをとってこれる
	// rulesを追加することでruleにそったものだけをストリームできる(ルールごと512制限があるので、同時に20ユーザぐらいのTweetをストリームできそう)
	// https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/api-reference/post-tweets-search-stream-rules
	e.Start(":8000")
}
