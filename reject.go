package xerr

type reject struct {
	Response []byte
	ShouldRecord bool
}
func (reject reject) Error() string {
	return string(reject.Response)
}
func ErrorToReject(err error) (rejectValue *reject, isReject bool) {
	switch err.(type) {
	case *reject:
		return err.(*reject), true
	default:
		return &reject{}, false
	}
}
func NewReject(response []byte, shouldRecord bool) error {
	return &reject{
		Response: response,
		ShouldRecord: shouldRecord,
	}
}
