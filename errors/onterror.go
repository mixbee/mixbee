

package errors

type mixbeeError struct {
	errmsg    string
	callstack *CallStack
	root      error
	code      ErrCode
}

func (e mixbeeError) Error() string {
	return e.errmsg
}

func (e mixbeeError) GetErrCode() ErrCode {
	return e.code
}

func (e mixbeeError) GetRoot() error {
	return e.root
}

func (e mixbeeError) GetCallStack() *CallStack {
	return e.callstack
}
