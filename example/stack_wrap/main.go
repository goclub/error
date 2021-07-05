package main

import (
	"database/sql"
	xerr "github.com/goclub/error"
	"log"
	"os"
)

func main () {
	log.Print("xerr.New")
	{
		err := xerr.New("abc")
		log.Print(err) // abc
	}
	log.Print("xerr.WrapPrefix")
	{
		err := biz()
		// 简单错误 biz: db: sql: transaction has already been committed or rolled back
		log.Print(err)

		// 附带堆栈信息 biz: db: sql: transaction has already been committed or rolled back
		// 等同于 log.Printf("%+v", err)
		xerr.LogWithStack(err)

		// 建议使用 xerr.Is 比较
		log.Print(xerr.Unwrap(err) == sql.ErrTxDone) // true
		log.Print(xerr.Is(err, sql.ErrTxDone)) // true
	}
}

func biz() error {
	err := db()
	return xerr.WrapPrefix("biz", err)
}
func db() error {
	return xerr.WrapPrefix("db", sql.ErrTxDone)
}

func openFile() error {
	_, err := os.Open("notexistfile")
	return err
}