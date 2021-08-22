package xerr_test

import (
	"errors"
	xerr "github.com/goclub/error"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIs(t *testing.T) {
	err := f1()
	assert.Equal(t,true, xerr.Is(err, ErrSome))
	assert.Equal(t,true,xerr.Is(ErrSome, err))
}
func f1() error {
	return f2()
}
func f2() error {
	return f3()
}
var ErrSome = errors.New("some")
func f3() error {
	return xerr.WithStack(ErrSome)
}
