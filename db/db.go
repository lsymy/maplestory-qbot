package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() error {
	var err error

	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/maple-story")

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println("db init success !")
		return nil
	}
}

func GetDB() *sql.DB {
	return db
}
