# error


## error  和 panic

> 只要执行 panic 就极大可能导致程序中断进程退出  

在 Web Server 服务中只有启动程序时才应该出现 panic 代码。

例如启动时错误

[panic_on_startup|embed](./docs/panic_on_startup/main.go)

在 main 函数中 如果出现数据库连接错误是无法处理的所以使用panic。

> 在某些场景下为了解决数据库偶尔连接失败但立即会恢复正常的情况，会在 db.Ping() 错误时不panic,而是记录日志报警。

例如某个URL 解析参数错误

[http_query_parse_error|embed](./docs/http_query_parse_error/main.go)

如果在获取请求参数并转换为数字时错误就 panic 是没有必要的。
因为 web 是面向很多用户的，不能因为某个接口被使用者传递了错误的参数就使用 panic 中断服务。
虽然 go http 路由一般都会在函数 panic 时候进行处理防止服务中断，但是 goroutine panic 如果没有 defer recover 会导致进程退出。

应该记住 **只要执行 panic 就有极大可能中断程序**。  

