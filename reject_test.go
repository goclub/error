package xerr_test

import (
	xerr "github.com/goclub/error"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
			return xerr.Reject(1, "abc", false)
		}()
		reject, isReject := xerr.AsReject(err)
		assert.Equal(t,reject.Message, "abc")
		assert.Equal(t,isReject, true)
		assert.Equal(t,reject.Code, int32(1))

	}
	{
		err := func () error {
			return xerr.Reject(1, "abc", true)
		}()
		reject, isReject := xerr.AsReject(err)
		assert.Equal(t,reject.Message, "abc")
		assert.Equal(t,isReject, true)
		assert.Equal(t,reject.Code, int32(1))
	}
	{
		err := func () error {
			return xerr.RejectWithPrivateDetails(1, "abc", xerr.PrivateDetails{"privateDetails"})
		}()
		reject, isReject := xerr.AsReject(err)
		assert.Equal(t,reject.Message, "abc")
		assert.Equal(t,isReject, true)
		assert.Equal(t,reject.Code, int32(1))
		assert.Equal(t,reject.PrivateDetails(), "privateDetails")
	}
}

func a() error {
	return b()
}
func b () error {
	return xerr.Reject(1, "abc", false)
}
