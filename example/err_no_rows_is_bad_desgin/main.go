package main

import (
	"database/sql"
	"log"
)

func BadDesign() {
	log.Print("BadDesign")
	db := &sql.DB{} // 为了演示直接获取结构体，实际上必须通过 sql.Open 才能拿到 sql.DB
	var name string
	row := db.QueryRow("select name from user where id= ?", 1)
	err := row.Scan(&name)
	if err != nil {
		// 数据是否存在要通过 err == sq.ErrNoRows 判断非常繁琐导致代码容易出错,而且就算是判断也应该用 xerr.Is
		if err == sql.ErrNoRows {
			log.Print("没数据")
			return
		} else {
			log.Print("err")
			return
		}
	}
	log.Print("name", name)
}
func GoodsDesgin() {
	// 如果 row.Scan() 设计成多返回一个 bool 来表示数据是否存在会更好
	// sq.ErrNoRows 虽然是官方库的设计，但我认为这是一种错误使用 error 的例子
	/*
	row := db.QueryRow("select name from user where id= ?", 1)
	var has bool
	has, err := row.Scan(&name)
	if err != nil {
		log.Print(err) ; return
	}
	if has == false {
		log.Print("没有数据") ; return
	}
	log.Print("name", name)
	*/
}
func main() {
	log.Print("注意这段代码因为没有连接数据库所以是无法运行的,只是为了演示错误的设计和好的设计")
	BadDesign()
	GoodsDesgin()
}

