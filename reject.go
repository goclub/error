package xerr

import (
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"testing"
)
type Resp struct {
	Error RespError `json:"error"`
}
type RespError struct {
	Code int32 `json:"code"`
	Message string `json:"message"`
}
type reject struct {
	Code int32
	Message string
	ShouldRecord bool
	*stack
}
func (reject reject) Error() string {
	return strconv.Itoa(int(reject.Code))+ ":" + reject.Message
}
func AsReject(err error) (rejectValue *reject, asReject bool) {
	asReject = As(err, &rejectValue)
	return
}
func Reject(code int32, message string, shouldRecord bool) error {
	if code == 0 {
		log.Print("xerr.Reject(code, message, shouldRecord) code can not be zero")
		code = 1
	}
	return &reject{
		Code: code,
		Message: message,
		ShouldRecord: shouldRecord,
		stack: callers(),
	}
}
func EqualRejectCode(t *testing.T, err error, code int32) {
	reject, asReject := AsReject(err)
	assert.True(t, asReject, err.Error(), " not xerr.reject")
	if asReject {
		assert.Equal(t, reject.Code, code)
	}
}
