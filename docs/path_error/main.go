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


