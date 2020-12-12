package xerr_test

import (
	xerr "github.com/goclub/error"
	"log"
	"strconv"
	"testing"
)

func TestMustInt(t *testing.T) {
	{
		i, err := strconv.Atoi("1") ; if err != nil {panic(err)}
		log.Print(i)
	}
	// 等同于
	{
		i := xerr.MustInt(strconv.Atoi("1"))
		log.Print(i)
	}
}