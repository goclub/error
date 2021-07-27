# goclub/error

```sh
go get github.com/goclub/error
```

## 堆栈和包装

官方的 "errors" 包没有堆栈,项目中用起来调试会很困难.

这是个大坑,如果要黑 go 一定要黑 go 的错误堆栈,黑泛型和 err != nil 太低级.

解决方案代表性的有2个

1. pkg/error https://github.com/pkg/error
2. xerror https://go.googlesource.com/proposal/+/master/design/29934-error-values.md


在go2没推出之前,还建议是使用 pkg/error 
但是 pkg/errors 的 package name 是 errors ,这导致经常与标准库 的 xerr.混淆.
为了避免混淆 xerr 调用了pkg/error 部分方法,去掉了一些以后会不兼容或多余的方法.

源码: [./error.go?embed](./error.go)

## error  和 panic

> 只要执行 panic 就极大可能导致程序中断进程退出  

在 Web Server 服务中只有启动程序时才应该出现 panic 代码。

例如启动时错误

[panic_on_startup](./example/panic_on_startup/main.go?embed)

在 main 函数中 如果出现数据库连接错误是无法处理的所以使用panic。

> 在某些场景下为了解决数据库偶尔连接失败但立即会恢复正常的情况，会在 db.Ping() 错误时不panic,而是记录日志报警。

例如某个URL 解析参数错误

[http_query_parse_error](./example/http_query_parse_error/main.go?embed)

如果在获取请求参数并转换为数字时错误就 panic 是没有必要的。
因为 web 是面向很多用户的，不能因为某个接口被使用者传递了错误的参数就使用 panic 中断服务。
虽然 go http 路由一般都会在函数 panic 时候进行处理防止服务中断，但是 goroutine panic 如果没有 defer recover 会导致进程退出。

应该记住 **只要执行 panic 就有极大可能中断程序**。
   

## sql.ErrNoRows 糟粕 

应当 避免 sql.ErrNoRows 这种错误的设计

[err_no_rows_is_bad_desgin](./example/err_no_rows_is_bad_desgin/main.go?embed)

数据不存在应该明确的通过 bool 表达，而不是让使用者通过 `err == sq.ErrNoRows` 来判断。

在日常开发中我们也应该避免写出类似的滥用 error 代码。

包装后直接使用 `err == pakcgaeName.ErrSome` 判断会失败，可以通过 `xerr.Is()` 解决。

## Sentinel Error

Sentinel 是哨兵的意思,哨兵这个名字取得很难理解.只有看过示例代码才知道是什么意思.

```go
// Sentinel Error
if err != nil {
	if xerr.Is(packageName.ErrName) {
      // do some	
    } else {
    	return err
    }
}
```   

通过比对 Sentinel Error 的方式判断需要借助文档才能弄清楚有哪些错误。

## 使用自定义错误类型携带更多的信息

os标准库有很多自定义错误类型的用法： 

[path_error](./example/path_error/main.go?embed)

判断错误类型的方式的缺点是不够直观，要基于约定和文档才能知道该如何判断错误。但这不妨碍在某些场景下使用自定义错误类型。

为了解决自定义类型被 `fmt.Errorf` 后类型断言不准确的问题，使用`xerr.As()` 进行判断


## reject

在日常的开发中有很多业务逻辑信息需要传递给客户端，例如创建用户时手机号码已存在。
如果逻辑层函数返回 `xerr.New("手机号码已存在")` 给协议层（http, rpc ）虽然能实现但是传递的信息太少，并且不安全。
因为协议层无法判断当前的错误是业务逻辑信息还是其他io错误.

使用 xerr.Reject() 创建可公开给用户的业务错误,使用 xerr.New 定义不公开的内部错误

[reject](./example/reject/main.go?embed)


源码实现非常简单,感兴趣可以看看

[源码](./reject.go?embed)

## defer


使用 defer 时需要注意不要赋值给 err 
```go
package main

import (
	"github.com/goclub/error"
	"log"
	"net/http"
)

func main () {
	log.Print(BugCode(0))
	log.Print(CorrectCode(0))
}
func BugCode(i int) (err error) {
	resp, err := http.Get("https://bing.com") ; if err != nil {
		return
	}
	defer func () {
		// 这里有个陷阱 Some(1)   BugCode(0) 永远不会返回  "i can not be zero"
		// 因为 err = 重新赋值了
		err = resp.Body.Close() ; if err != nil {
			return
		}
	}()
	if i == 0 {
		return xerr.New("i can not be zero")
	}
	return nil
}

func CorrectCode(i int) (err error) {
	resp, err := http.Get("https://bing.com") ; if err != nil {
		return
	}
	defer func () {
		// 不要覆 err
		closeErr := resp.Body.Close() ; if closeErr != nil {
			return
		}
	}()
	if i == 0 {
		return xerr.New("i can not be zero")
	}
	return nil
}
```

## 最佳实践

1. 除了启动(main)或者初始化(init)代码不要使用 panic
2. 避免出现 sql.ErrNoRows 这种滥用错误的设计
3. 使用 xerr.New 和 xerr.Reject区分服务内部错误和公开业务错误,避免泄露敏感信息.
4. 不知道如何处理的错误时向上传递
5. error 应该是个黑盒，大部分情况下使用者只需要判断 err != nil 进行错误处理或向上传递错误。
