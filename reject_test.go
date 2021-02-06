package xerr_test

import (
	"encoding/json"
	"fmt"
	xerr "github.com/goclub/error"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReject_Error(t *testing.T) {
	data, err := json.Marshal(map[string]string{"type":"pass"}) ; if err != nil {
		panic(err)
	}
	assert.Equal(t,xerr.NewReject(data, true).Error(), `{"type":"pass"}`)
	testInterface := func(err error) {/* 编译期不报错即可 */}
	testInterface(xerr.NewReject(nil, false))
}
func TestAsReject(t *testing.T) {
	{
		var err error
		err = nil
		reject, isReject := xerr.AsReject(err)
		assert.Nil(t, reject)
		assert.Equal(t,isReject, false)
	}
	{
		err := func () error {
			return xerr.NewReject([]byte("abc"), false)
		}()
		reject, isReject := xerr.AsReject(err)
		assert.Equal(t,reject, xerr.NewReject([]byte("abc"), false))
		assert.Equal(t,isReject, true)
	}
	{
		err := func () error {
			return xerr.NewReject([]byte("abc"), true)
		}()
		reject, isReject := xerr.AsReject(err)
		assert.Equal(t,reject, xerr.NewReject([]byte("abc"), true))
		assert.Equal(t,isReject, true)
	}

}
func Some() error {
	return xerr.NewReject(NewFail("用户不存在"), false)
	// return nil
}
func TestReject(t *testing.T) {
	{
		err := Some()
		assert.NotNil(t, err)
		reject, isReject := xerr.AsReject(err)
		assert.Equal(t, isReject, true)
		assert.Equal(t, reject, xerr.NewReject(NewFail("用户不存在"), false))
	}
	{
		err := a()
		assert.NotNil(t, err)
		reject, isReject := xerr.AsReject(err)
		assert.Equal(t, isReject, true)
		assert.Equal(t, reject, xerr.NewReject([]byte("b"), false))
		assert.Equal(t, err.Error(), "a: b")
	}
}
func a() error {
	err  := b()
	return fmt.Errorf("a: %w", err)
}

func b() error {
	return xerr.NewReject([]byte("b"), false)
}
type ResponseType string
func (v ResponseType) Enum() (enum struct{
	Pass ResponseType
	Fail ResponseType
}) {
	enum.Pass = "pass"
	enum.Fail = "fail"
	return
}
type Response struct {
	Type ResponseType `json:"type"`
	Msg string `json:"msg"`

}
func NewPass(data interface{}) Response {
	return Response{
		Type: Response{}.Type.Enum().Pass,
	}
}
func NewFail(msg string) []byte {
	data, err := json.Marshal(Response{
		Type: Response{}.Type.Enum().Fail,
		Msg: msg,
	}) ; if err != nil {
		panic(err)
	}
	return  data
}

func QueryRow(sql string, v interface{}) error {
	return nil
}


