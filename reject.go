package xerr

import (
	"errors"
	"strconv"
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
func NewReject(code int32, message string, shouldRecord bool) error {
	return &reject{
		Code: code,
		Message: message,
		ShouldRecord: shouldRecord,
	}
}
