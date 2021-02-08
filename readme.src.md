# goclub/error


## error  和 panic

> 只要执行 panic 就极大可能导致程序中断进程退出

在 Web Server 服务中只有启动程序时才应该出现 panic 代码。

例如启动时错误

[panic_on_startup|embed](./docs/panic_on_startup/main.go)

在 main 函数中 如果出现数据库连接错误是无法处理的所以使用 `panic`。

> 在某些场景下为了解决数据库偶尔连接失败但立即会恢复正常的情况，会在 db.Ping() 错误时不panic,而是记录日志报警。

---

非启动代码则不会使用 panic

例如某个 http 接口解析参数时错误

[http_query_parse_error|embed](./docs/http_query_parse_error/main.go)

如果在获取请求参数并转换为数字时错误就 panic 是没有必要的。
因为 web 是面向很多用户的，不能因为某个接口被使用者传递了错误的参数就使用 panic 中断服务。
虽然 go http 路由一般都会在函数 panic 时候进行处理防止服务中断，但是 goroutine panic 如果没有 defer recover 会导致进程退出。

应该记住 **只要执行 panic 就有极大可能中断程序**。
   
## 避免 sql.ErrNoRows 这种错误的设计

[err_no_rows_is_bad_desgin|embed](./docs/err_no_rows_is_bad_desgin/main.go)

数据不存在应该明确的通过 bool 表达，而不是让使用者通过 `err == sq.ErrNoRows` 来判断。

在日常开发中我们也应该避免写出类似的滥用 error 代码。

## fmt.Errorf

可以通过 `fmt.Errorf` 包装错误后传递给上层，附带更多信息。

```go
err := some() if err !+ nil {
	return fmt.Errorf("abc: %w", err)
}
```

包装后直接使用 `err == pakcgaeName.ErrSome` 判断会失败，可以通过 `errors.Unwrap()` `errors.Is())` 解决。

## Sentinel Error

```go
// Sentinel Error
if err != nil {
	if errors.Is(packageName.ErrName) {
      // do some	
    } else {
    	return err
    }
}
```   

通过比对 Sentinel Error 的方式判断需要借助文档才能弄清楚有哪些错误。

## 使用自定义错误类型携带更多的信息

os标准库有很多自定义错误类型的用法： [path_error|embed](./docs/path_error/main.go)

判断错误类型的方式的缺点是不够直观，要基于约定和文档才能知道该如何判断错误。但这不妨碍在某些场景下使用自定义错误类型。

为了解决自定义类型被 `fmt.Errorf` 后类型断言不准确的问题，使用`errors.As()` 进行判断

## reject

在日常的开发中有很多业务逻辑信息需要传递给客户端，例如创建用户时手机号码已存在。
如果逻辑层函数返回 `errors.New("手机号码已存在")` 给协议层（http, rpc ）虽然能实现但是传递的信息太少，并且不安全。
因为协议层无法判断当前的错误是业务逻辑信息还是其他io错误，例如 :

[reject|embed](./docs/reject/main.go)

源码实现非常简单,感兴趣可以看看

[源码|embed](./reject.go)

## 最佳实践

1. 除了启动(main)或者初始化(init)代码不要使用 panic
2. 避免出现 sql.ErrNoRows 这种滥用错误的设计
3. 若一定要包含错误信息则参考 xerr.NewReject 实现自定义错误类型
4. 不知道如何处理的错误时向上传递
5. error 应该是个黑盒，大部分情况下使用者只需要判断 err != nil 进行错误处理或向上传递错误。
