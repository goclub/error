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
			return xerr.NewReject(0, "abc", false)
		}()
		reject, isReject := xerr.AsReject(err)
		assert.Equal(t,reject, xerr.NewReject(0, "abc", false))
		assert.Equal(t,isReject, true)
		xerr.TestEqualCode(t, err, 0)

	}
	{
		err := func () error {
			return xerr.NewReject(0, "abc", true)
		}()
		reject, isReject := xerr.AsReject(err)
		assert.Equal(t,reject, xerr.NewReject(0, "abc", true))
		assert.Equal(t,isReject, true)
	}

}

