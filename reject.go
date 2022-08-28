package xerr

import (
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"testing"
)

type IRejecter interface {
	Error() string
	Resp() Resp
}

// reject 被设计成不公开类型,如果需要在自己的项目选中传递 reject,可以使用 IRejecter 接口传递
type reject struct {
	Code         int32  `json:"-" xml:"-"`
	Message      string `json:"-" xml:"-"`
	ShouldRecord bool   `json:"-" xml:"-"`
}

func (reject reject) Error() string {
	return strconv.Itoa(int(reject.Code)) + ":" + reject.Message
}
func (reject reject) Resp() Resp {
	return NewResp(reject.Code, reject.Message)
}
func AsReject(err error) (rejectValue *reject, asReject bool) {
	asReject = As(err, &rejectValue)
	return
}
func Reject(code int32, publicMessage string, shouldRecord bool) error {
	if code == 0 {
		log.Print("xerr.Reject(code, message, shouldRecord) code can not be zero")
		code = 1
	}
	return WithStack(&reject{
		Code:         code,
		Message:      publicMessage,
		ShouldRecord: shouldRecord,
	})
}
func EqualRejectCode(t *testing.T, err error, code int32) {
	reject, asReject := AsReject(err)
	assert.True(t, asReject, err.Error(), " not xerr.reject")
	if asReject {
		assert.Equal(t, reject.Code, code)
	}
}
func EqualRejectMessage(t *testing.T, err error, message string) {
	reject, asReject := AsReject(err)
	assert.True(t, asReject, err.Error(), " not xerr.reject")
	if asReject {
		assert.Equal(t, reject.Message, message)
	}
}
