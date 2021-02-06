package xerr

import "errors"

type reject struct {
	Response []byte
	ShouldRecord bool
}
func (reject *reject) Error() string {
	return string(reject.Response)
}
func AsReject(err error) (rejectValue *reject, isReject bool) {
	isReject = errors.As(err, &rejectValue)
	return
}
func NewReject(response []byte, shouldRecord bool) error {
	return &reject{
		Response: response,
		ShouldRecord: shouldRecord,
	}
}
