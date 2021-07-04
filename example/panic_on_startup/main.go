package  main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goclub/sql"
)

func main () {
	db, dbClose, err := sq.Open("mysql", sq.MysqlDataSource{
		User: "root",
		Password:"",
		Host: "127.0.0.1",
		Port:"3306",
		DB: "test_goclub_sql",
	}.FormatDSN()) ; if err != nil {
		// 无法处理的错误
		panic(err)
	}
	defer dbClose()
	err = db.Ping(context.TODO()) ; if err != nil {
		// 无法处理的错误
		panic(err)
	}
}