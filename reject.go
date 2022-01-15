package xerr

import (
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"testing"
)
type reject struct {
	Code int32
	Message string
	ShouldRecord bool
	// 设置私有字段防止被意外 json marshal 导致泄漏信息
	privateDetails string
}
func (reject reject) Error() string {
	return strconv.Itoa(int(reject.Code))+ ":" + reject.Message
}
func (reject reject) Resp() Resp {
	return NewResp(reject.Code, reject.Message)
}
func (reject reject) PrivateDetails() string {
	return reject.privateDetails
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
		Code: code,
		Message: publicMessage,
		ShouldRecord: shouldRecord,
	})
}
type PrivateDetails struct {PrivateDetail string}
func RejectWithPrivateDetails(code int32, publicMessage string, privateDetails PrivateDetails) (err error) {
	if code == 0 {
		log.Print("xerr.Reject(code, message, shouldRecord) code can not be zero")
		code = 1
	}
	return WithStack(&reject{
		Code: code,
		Message: publicMessage,
		ShouldRecord: true,
		privateDetails: privateDetails.PrivateDetail,
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
