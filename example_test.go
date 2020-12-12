package xerr_test

import (
	_ "github.com/go-sql-driver/mysql"
	xerr "github.com/goclub/error"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
)

// func TestExample(t *testing.T) {
// 	go ExampleHttpServer()
// 	go ExampleRejectHttpServer()
// 	select{} // select{} 用于让主routine一直等待
// }

// 仅作为演示使用 map[string]bool 存储数据，请暂时忽略没加锁导致的并发竞态问题。
var users = map[string]bool{}
func Register(name string)  error {
	if len(name) == 0 {
		return errors.New("name 必填")
	}
	_, hasUser := users[name]
	if hasUser {
		return errors.New("用户" + name + "已存在")
	} else {
		users[name] =  true
	}
	return nil
}
func NewsCreate(title string) error {
	if len(title) == 0 {
		return errors.New("title必填")
	}
		_, err := http.PostForm("http://www.someerrorapi.com/news_create?ak=abc&sk=password", url.Values{"title":[]string{title}}) ; if err != nil {
		return err
	}
	return nil
}
func WriteMessage(writer http.ResponseWriter, data []byte) {
	_, writeErr := writer.Write(data) ; if writeErr != nil {panic(writeErr)}
}
// 一个场景的将错误传递到协议层（http rpc）的示例
// http://127.0.0.1:1341/register 通过 error 传递了一些供用户理解的提示信息
// http://127.0.0.1:1341/news/create 接口调用失败时会将 url 中的 ak sk 等安全信息保留，这会造成严重的安全问题
// 使用 xerr.Reject 区分请求拒绝信息和内部错误可以避免这类安全问题
func ExampleHttpServer() {
	URLRegister := "/register"
	URLNewsCreate := "/news/create"
	http.HandleFunc(URLRegister, func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		err := Register(query.Get("name")) ; if err != nil {
			WriteMessage(writer, []byte(err.Error())) ; return
		}
		WriteMessage(writer, []byte("register success")) ; return
	})
	http.HandleFunc(URLNewsCreate, func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		err := NewsCreate(query.Get("title")) ; if err != nil {
			WriteMessage(writer, []byte(err.Error())) ; return
		}
		WriteMessage(writer, []byte("create news success")) ; return
	})
	addr := ":1341"
	log.Print("http://127.0.0.1" + addr + URLRegister)
	log.Print("http://127.0.0.1" + addr + URLNewsCreate)
	log.Print(http.ListenAndServe(addr, nil))
}

func ArticleCreate(title string) error {
	if len(title) == 0 {
		// 一些正常情况下不应该出现的错误，NewReject() 第二个参数为 true
		// 例如用户编辑一个不属于自己的文章时使用 return xerr.NewReject("你没有次文字的访问权限", true)
		return xerr.NewReject([]byte("title必填"), false)
	}
	_, err := http.PostForm("http://www.someerrorapi.com/article_create?ak=abc&sk=password", url.Values{"title":[]string{title}}) ; if err != nil {
		return err
	}
	return nil
}

func ExampleRejectHttpServer() {
	URLArticleCreate := "/article/create"
	http.HandleFunc(URLArticleCreate, func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		err := ArticleCreate(query.Get("title")) ; if err != nil {
			reject, isReject := xerr.ErrorToReject(err)
			if isReject {
				if reject.ShouldRecord {log.Print(reject.Response)}
				WriteMessage(writer, reject.Response) ; return
			} else {
				writer.WriteHeader(http.StatusInternalServerError)
				WriteMessage(writer, []byte("服务出错")) ; return
			}
		}
		WriteMessage(writer, []byte("create article success")) ; return
	})
	addr := ":6341"
	log.Print("http://127.0.0.1" + addr + URLArticleCreate)
	log.Print("http://127.0.0.1" + addr + URLArticleCreate +"?title=abc")
	log.Print(http.ListenAndServe(addr, nil))
}