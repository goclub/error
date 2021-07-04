package xerr

import (
	"fmt"
	"golang.org/x/xerrors"
)

func New(text string) error {
	return xerrors.New(text)
}

func Is(err, target error) bool {
	return xerrors.Is(err, target)
}
func Unwrap(err error) error {
	return xerrors.Unwrap(err)
}
func WrapPrefix(prefix string, err error) error {
	return fmt.Errorf(prefix + "%w", err)

}
func Errorf(format string, a ...interface{}) error {
	return xerrors.Errorf(format, a...)
}
func As(err error, target interface{}) bool {
	return xerrors.As(err, target)
}
type Frame xerrors.Frame
func Caller(skip int) Frame{
	return Frame(xerrors.Caller(skip))
}
type Formatter interface {
	error
	FormatError(p xerrors.Printer) (next error)
}

func FormatError(f Formatter, s fmt.State, verb rune) {
	xerrors.FormatError(f, s, verb)
}
func Opaque(err error) error {
	return xerrors.Opaque(err)
}
