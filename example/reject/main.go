package main

import (
	"context"
	"fmt"
	xerr "github.com/goclub/error"
	xhttp "github.com/goclub/http"
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
	url := "/create_user?ak=" + ak +"&sk=" + sk
	origin := "http://www.exist-domain-only-test.com"
	if name == "admin" {
		return xerr.New("name can not be admin")
	}
	_, bodyClose, statusCode, err := client.Send(context.TODO(), "GET", origin, url, xhttp.SendRequest{}) ; if err != nil {
		return err
	}
	defer bodyClose()
	if statusCode != 200 {
		return xerr.New("statusCode error:" + strconv.Itoa(statusCode))
	}
	return nil
}
func CreateUserSafeReject(name string) error {
	if name == "admin" {
		return xerr.Reject(1, "name can not be admin", false)
	}
	ak := "nimoc"
	sk := "1234"
	url := "/create_user?ak=" + ak +"&sk=" + sk
	origin := "http://www.exist-domain-only-test.com"
	_, bodyClose, statusCode, err := client.Send(context.TODO(), "GET", origin, url, xhttp.SendRequest{}) ; if err != nil {
		return fmt.Errorf("create user:%w, err")
	}
	defer bodyClose()
	if statusCode != 200 {
		return xerr.New("statusCode error:" + strconv.Itoa(statusCode))
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
		// 你实现一些拦截器来避免重复代码.或者直接使用 github.com/goclub/http
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
