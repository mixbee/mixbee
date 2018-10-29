

package errors

import (
	"errors"
)

const callStackDepth = 10

type DetailError interface {
	error
	ErrCoder
	CallStacker
	GetRoot() error
}

func NewErr(errmsg string) error {
	return errors.New(errmsg)
}

func NewDetailErr(err error, errcode ErrCode, errmsg string) DetailError {
	if err == nil {
		return nil
	}

	onterr, ok := err.(mixbeeError)
	if !ok {
		onterr.root = err
		onterr.errmsg = err.Error()
		onterr.callstack = getCallStack(0, callStackDepth)
		onterr.code = errcode

	}
	if errmsg != "" {
		onterr.errmsg = errmsg + ": " + onterr.errmsg
	}

	return onterr
}

func RootErr(err error) error {
	if err, ok := err.(DetailError); ok {
		return err.GetRoot()
	}
	return err
}
