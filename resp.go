package xerr

type Resp struct {
	Error respError `json:"error"`
}
func NewResp(code int32, message string) Resp {
	return Resp{
		Error: respError{
			Code:    code,
			Message: message,
		},
	}
}
type respError struct {
	Code int32 `json:"code"`
	Message string `json:"message"`
}