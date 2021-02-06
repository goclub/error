package  main

import (
	"github.com/goclub/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main () {
	db, dbClose, err := sq.Open("mysql", sq.DataSourceName{
		DriverName: "mysql",
		User: "root",
		Password:"somepass",
		Host: "127.0.0.1",
		Port:"3306",
		DB: "test_goclub_sql",
	}.String()) ; if err != nil {
		// 无法处理的错误
		panic(err)
	}
	defer dbClose()
	err = db.Ping() ; if err != nil {
		// 无法处理的错误
		panic(err)
	}
}