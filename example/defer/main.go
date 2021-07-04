package main

import (
	xerr "github.com/goclub/error"
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