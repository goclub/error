package xerr

import (
	pkgErr "github.com/pkg/errors"
	"log"
)

var (
	// 创建错误
	New    = pkgErr.New
	// 通过format 创建错误
	Errorf = pkgErr.Errorf
	// 判断自定义错误
	As     = pkgErr.As
	WithStack = pkgErr.WithStack
)
// 包装错误
func WrapPrefix(prefix string, err error) error {
	return pkgErr.Wrap(err, prefix)
}
// 判断 Sentinel Error 错误
func Is(err error, target error) bool {
	if pkgErr.Is(err, target) {
		return true
	}
	return pkgErr.Is(target, err)
}

// 获取被包装的底层错误
func Unwrap(err error) error{
	e := pkgErr.Cause(err)
	if e != nil {
		return e
	}
	return pkgErr.Unwrap(err)
}
func PrintStack(err error) {
	log.Printf("%+v", err)
}
