package xerr

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

type reject struct {
	Code int32
	Message string
	ShouldRecord bool
}
func (reject reject) Error() string {
	return strconv.Itoa(int(reject.Code))+ ":" + reject.Message
}
func AsReject(err error) (rejectValue *reject, asReject bool) {
	asReject = errors.As(err, &rejectValue)
	return
}
func Reject(code int32, message string, shouldRecord bool) error {
	return &reject{
		Code: code,
		Message: message,
		ShouldRecord: shouldRecord,
	}
}
func TestEqualCode(t *testing.T, err error, code int32) {
	reject, asReject := AsReject(err)
	assert.True(t, asReject, err.Error(), " not xerr.reject")
	if asReject {
		assert.Equal(t, reject.Code, code)
	}
}
// func TestEqualMessage(t *testing.T, err error, message string) {
// 	reject, asReject := AsReject(err)
// 	assert.True(t, asReject, err.Error(), " not xerr.reject")
// 	if asReject {
// 		assert.Equal(t, reject.Message, message)
// 	}
// }
// func TestEqualShouldRecord(t *testing.T, err error, shouldRecord bool) {
// 	reject, asReject := AsReject(err)
// 	assert.True(t, asReject, err.Error(), " not xerr.reject")
// 	if asReject {
// 		assert.Equal(t, reject.Message, shouldRecord)
// 	}
// }