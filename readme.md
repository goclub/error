# error


## error  和 panic

> 只要执行 panic 就极大可能导致程序中断进程退出  

在 Web Server 服务中只有启动程序时才应该出现 panic 代码。

例如启动时错误

[panic_on_startup](./docs/panic_on_startup/main.go)
```.go
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
```

在 main 函数中 如果出现数据库连接错误是无法处理的所以使用panic。

> 在某些场景下为了解决数据库偶尔连接失败但立即会恢复正常的情况，会在 db.Ping() 错误时不panic,而是记录日志报警。

例如某个URL 解析参数错误

[http_query_parse_error](./docs/http_query_parse_error/main.go)
```.go
package main

import (
	"log"
	"net/http"
	"strconv"
)

func main () {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request){
		query := request.URL.Query()

		page, err := strconv.Atoi(query.Get("page")) ; if err != nil {
			writer.WriteHeader(500)
			log.Print(err)
		}

		log.Print("page", page)

		writer.WriteHeader(200)
		_, err = writer.Write([]byte("ok")) ; if err != nil {
			writer.WriteHeader(500)
			log.Print(err)
		}
	})
	err := http.ListenAndServe(":3000", mux) ; if err != nil {
		// 服务无法启动必须 panic
		panic(err)
	}
}

```

如果在获取请求参数并转换为数字时错误就 panic 是没有必要的。
因为 web 是面向很多用户的，不能因为某个接口被使用者传递了错误的参数就使用 panic 中断服务。
虽然 go http 路由一般都会在函数 panic 时候进行处理防止服务中断，但是 goroutine panic 如果没有 defer recover 会导致进程退出。

应该记住 **只要执行 panic 就有极大可能中断程序**。
   

## 避免 sql.ErrNoRows 这种错误的设计

[err_no_rows_is_bad_desgin](./docs/err_no_rows_is_bad_desgin/main.go)
```.go
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
		// 数据是否存在要通过 err == sq.ErrNoRows 判断非常繁琐导致代码容易出错
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
	BadDesign()
	GoodsDesgin()
}


```

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

os标准库有很多自定义错误类型的用法： [path_error](./docs/path_error/main.go)
```.go
package main

import (
	"errors"
	"log"
	"os"
)

func main() {
	f, err := os.Open("/nonexistent.txt")
	// err = fmt.Errorf("some: %w", err) // 因为下面的代码使用了 errors.As 判断，即使错误被 Errorf 包装了依然能通过 Unwrap() 判断原始错误是否一致。
	if err != nil {
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			log.Print(pathError.Op, pathError.Path, " failed")
			return
		} else {
			log.Print(err)
		}
	}
	log.Print(f.Name(), "opened successfully")
}



```

判断错误类型的方式的缺点是不够直观，要基于约定和文档才能知道该如何判断错误。但这不妨碍在某些场景下使用自定义错误类型。

为了解决自定义类型被 `fmt.Errorf` 后类型断言不准确的问题，使用`errors.As()` 进行判断


## reject

在日常的开发中有很多业务逻辑信息需要传递给客户端，例如创建用户时手机号码已存在。
如果逻辑层函数返回 `errors.New("手机号码已存在")` 给协议层（http, rpc ）虽然能实现但是传递的信息太少，并且不安全。
因为协议层无法判断当前的错误是业务逻辑信息还是其他io错误，例如 :

[reject](./docs/reject/main.go)
```.go
package main

import (
	"context"
	"fmt"
	xerr "github.com/goclub/error"
	xhttp "github.com/goclub/http"
	"errors"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
)

func WriteBytes(w http.ResponseWriter, data []byte) {
	_, err :=w.Write(data) ; if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}
}
var client = xhttp.NewClient(&http.Client{})
func CreateUserUnsafeError(name string) error {
	ak := "nimoc"
	sk := "1234"
	url := "http://www.exist-domain-only-test.com/create_user?ak=" + ak +"&sk=" + sk
	if name == "admin" {
		return errors.New("name can not be admin")
	}
	_, bodyClose, statusCode, err := client.Send(context.TODO(), "GET", url, xhttp.SendRequest{}) ; if err != nil {
		return err
	}
	defer bodyClose()
	if statusCode != 200 {
		return errors.New("statusCode error:" + strconv.Itoa(statusCode))
	}
	return nil
}
func CreateUserSafeReject(name string) error {
	if name == "admin" {
		// 没有 code 时 传 0
		return xerr.NewReject(0, "name can not be admin", false)
	}
	ak := "nimoc"
	sk := "1234"
	url := "http://www.exist-domain-only-test.com/create_user?ak=" + ak +"&sk=" + sk
	_, bodyClose, statusCode, err := client.Send(context.TODO(), "GET", url, xhttp.SendRequest{}) ; if err != nil {
		return fmt.Errorf("create user:%w, err")
	}
	defer bodyClose()
	if statusCode != 200 {
		return errors.New("statusCode error:" + strconv.Itoa(statusCode))
	}
	return nil
}
func main () {
	mux := http.NewServeMux()
	mux.HandleFunc("/unsafe_error", func(writer http.ResponseWriter, request *http.Request) {
		err := CreateUserUnsafeError(request.URL.Query().Get("name"))
		if err != nil {
			// 无论是什么错误都将 err.Error 响应给客户端
			WriteBytes(writer, []byte(err.Error())) ; return
		}
		WriteBytes(writer, []byte("ok"))
	})
	mux.HandleFunc("/safe_reject", func(writer http.ResponseWriter, request *http.Request) {
		err := CreateUserSafeReject(request.URL.Query().Get("name"))
		if err != nil {
			if reject, asReject := xerr.AsReject(err); asReject {
				if reject.ShouldRecord { log.Print(reject) }
				WriteBytes(writer, []byte("code("+ strconv.Itoa(int(reject.Code)) + ") " + reject.Message)) ; return
			} else {
				debug.PrintStack()
				log.Print(err)
				// writer.WriteHeader(500)
				WriteBytes(writer, []byte("server error")) ; return
			}
		}
		WriteBytes(writer, []byte("ok"))
	})
	addr := ":3000"
	log.Print("unsafe_error 会暴露 ak sk")
	log.Print("http://127.0.0.1" + addr + "/unsafe_error")
	log.Print("http://127.0.0.1" + addr + "/unsafe_error?name=admin")
	log.Print("safe_reject 不会暴露 ak sk")
	log.Print("http://127.0.0.1" + addr + "/safe_reject")
	log.Print("http://127.0.0.1" + addr + "/safe_reject?name=admin")
	log.Print(http.ListenAndServe(addr, mux))
}

```

源码实现非常简单,感兴趣可以看看

[源码](./reject.go)
```.go
package xerr

import "errors"

type reject struct {
	Response []byte
	ShouldRecord bool
}
func (reject *reject) Error() string {
	return string(reject.Response)
}
func AsReject(err error) (rejectValue *reject, isReject bool) {
	isReject = errors.As(err, &rejectValue)
	return
}
func NewReject(response []byte, shouldRecord bool) error {
	return &reject{
		Response: response,
		ShouldRecord: shouldRecord,
	}
}
```

## 最佳实践

1. 除了启动(main)或者初始化(init)代码不要使用 panic
2. 避免出现 sql.ErrNoRows 这种滥用错误的设计
3. 若一定要包含错误信息则参考 xerr.NewReject 实现自定义错误类型
4. 不知道如何处理的错误时向上传递
5. error 应该是个黑盒，大部分情况下使用者只需要判断 err != nil 进行错误处理或向上传递错误。
