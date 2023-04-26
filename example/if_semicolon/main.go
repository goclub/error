package main

import (
	"encoding/json"
	"log"
)

// # if err 陷阱

// 接下来会频繁出现 json.Unmarshal(nil, nil)
// 目的是为了获取 一个 errors.New("unexpected end of JSON input")

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

// ifSemicolon: 使用 if 分号 避免 err != nil 被意外移除
func ifSemicolon() (err error) {
	// gofmt 不会让 err != nil 与 Unmarshal 之间出现换行
	if err = json.Unmarshal(nil, nil); err != nil {
		return
	}
	return
}

// 自定义编辑器代码模版
// Jetbrains Goland live templates (支持 Surround With...)
// https://www.jetbrains.com/help/go/settings-live-templates.html
/*
	<template name="e" value="if $var$ = $SELECTION$; err != nil {&#10;    return&#10;}" description="if err = fn(); err != nil { return }" toReformat="false" toShortenFQNames="true">
	  <variable name="var" expression="errorVariableDefinition(SELECTION)" defaultValue="" alwaysStopAt="true" />
	  <context>
		<option name="GO" value="true" />
	  </context>
	</template>

	复制上面的代码 从 <template 到 </template> 进入偏好设置->编辑器->实时模版点击Go 粘贴即可创建
*/

func MultipleReturnValue() (err error) {
	var b []byte
	if b, err = json.Marshal(1); err != nil {
		return
	}
	log.Print("b:", string(b))
	return
	// if semicolon 的缺点是遇到多返回值时不能使用短变量声明 :=
	// 如果使用短变量声明:  if b, err := json.Marshal(1); err != nil {
	// b 就变为了 if 内部变量,并且 "if b," 处 会出现编译期报错 Unused variable 'b'
	// 但是瑕不掩瑜,不用短声明换取不出现"意外移除 err != nil"导致的bug是值得的
}

func main() {
	// 常见的代码
	// normal: unexpected end of JSON input
	log.Print("normal: ", normal())
	// 错误的代码
	// badCode: <nil>
	log.Print("badCode: ", badCode())
	// 解决方法: if semicolon
	// ifSemicolon: unexpected end of JSON input
	log.Print("ifSemicolon: ", ifSemicolon())
	// MultipleReturnValue: unexpected end of JSON input
	log.Print("MultipleReturnValue: ", MultipleReturnValue())

}
