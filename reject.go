package xerr

import (
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
	asReject = As(err, &rejectValue)
	return
}
func Reject(code int32, message string, shouldRecord bool) error {
	return &reject{
		Code: code,
		Message: message,
		ShouldRecord: shouldRecord,
	}
}
func EqualRejectCode(t *testing.T, err error, code int32) {
	reject, asReject := AsReject(err)
	assert.True(t, asReject, err.Error(), " not xerr.reject")
	if asReject {
		assert.Equal(t, reject.Code, code)
	}
}
