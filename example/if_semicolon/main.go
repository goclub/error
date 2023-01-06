package main

import (
	"encoding/json"
	"log"
)

// main: if err 陷阱
func main() {
	// 接下来会频繁出现 json.Unmarshal(nil, nil)
	// 目的是为了获取 一个 errors.New("unexpected end of JSON input")

	// 常见的代码
	// normal: unexpected end of JSON input
	log.Print("normal: ", normal())

	// 错误的代码
	// badCode: <nil>
	log.Print("badCode: ", badCode())

	// gofmt
	// 为了避免意外移除 err != nil 可以使用分号连接 函数与 err != nil
	// err = json.Unmarshal(nil, nil) ; if err != nil {
	// 	return
	// }
	// 但是使用gofmt格式化代码后会导致上面的被格式为
	// err = json.Unmarshal(nil, nil)
	// if err != nil {
	// 	return
	// }

	// 解决方法: if semicolon
	// ifSemicolon: unexpected end of JSON input
	log.Print("ifSemicolon: ", ifSemicolon())
	// MultipleReturnValue: unexpected end of JSON input
	log.Print("MultipleReturnValue: ", MultipleReturnValue())

	// 你可以在 IDE 中加入 live templates
	// 快捷键/缩写: e
	/* 模版内容:
	if err = $END$; err != nil {
			return
		}
	*/
}

// normal: 常见的错误返回
func normal() (err error) {
	err = json.Unmarshal(nil, nil) // 与下一行不可分隔
	if err != nil {                // 与上一行不可分隔
		return
	}
	return
}

// badCode: 错误的代码
func badCode() (err error) {
	err = json.Unmarshal(nil, nil) // 在修改迭代过程中 if err != nil 意外的被删除或移动
	var v uint
	// 此处赋值了 err 覆盖掉了上面 json.Unmarshal(nil, nil) 的 err
	err = json.Unmarshal([]byte(`1`), &v)
	if err != nil {
		return
	}
	// 最终 err 是 nil, json.Unmarshal(nil, nil) 失败了却没有返回错误
	return
}

// ifSemicolon: 使用 if 分号 避免 err != nil 被意外移除
func ifSemicolon() (err error) {
	// gofmt 不会让 err != nil 与 Unmarshal 之间出现换行
	if err = json.Unmarshal(nil, nil); err != nil {
		return
	}
	return
}

func MultipleReturnValue() (err error) {
	var b []byte
	if b, err = json.Marshal(1); err != nil {
		return
	}
	log.Print("b:", string(b))
	return
	// if semicolon 的缺点是遇到多返回值时不能使用短变量声明 :=
	// 如果使用短变量声明:  if b, err := json.Marshal(1); err != nil {
	// "if b," 处 会出现编译期报错 Unused variable 'b'
	// 瑕不掩瑜,不用短声明换取不出现"意外移除 err != nil"导致的bug是值得的
}
