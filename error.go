package xerr

import (
	pkgErr "github.com/pkg/errors"
)

var (
	// 创建错误
	New    = pkgErr.New
	// 通过format 创建错误
	Errorf = pkgErr.Errorf
	// 判断 Sentinel Error
	Is     = pkgErr.Is
	// 判断自定义错误
	As     = pkgErr.As
)
// 包装错误
func WrapPrefix(prefix string, err error) error {
	return pkgErr.Wrap(err, prefix)
}
// 获取被包装的底层错误
func Unwrap(err error) error{
	e := pkgErr.Cause(err)
	if e != nil {
		return e
	}
	return pkgErr.Unwrap(err)
}