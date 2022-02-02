package main

import (
	xerr "github.com/goclub/error"
	"log"
	"os"
)
func main() {
	f, err := os.Open("/nonexistent.txt")
	if err != nil {
		var pathError *os.PathError
		// 注意此处 &pathError 传递的是 os.PathError 指针的指针
		if xerr.As(err, &pathError) {
			log.Print(pathError.Op, pathError.Path, " failed")
			return
		} else {
			log.Print(err)
		}
	}
	log.Print(f.Name(), "opened successfully")
}


