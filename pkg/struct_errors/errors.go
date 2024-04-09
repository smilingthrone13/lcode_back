package struct_errors

import "github.com/pkg/errors"

type BaseError struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
	Err  error  `json:"-"`
}

func NewBaseErr(msg string, err error) *BaseError {
	e := &BaseError{}
	e.SetCode("default.base_error")
	//e.SetMsg(msg)
	e.SetErr(msg, err)

	return e
}

func (b *BaseError) SetCode(code string) {
	b.Code = code
}

func (b *BaseError) SetErr(msg string, e error) {
	b.Msg = msg

	var be *BaseError

	switch {
	case errors.As(e, &be):
		b.Err = be
	case e != nil:
		b.Err = &BaseError{
			Code: "default.base_error",
			Msg:  e.Error(),
			Err:  e,
		}
	default:
		b.Err = &BaseError{
			Code: "default.base_error",
			Msg:  b.Msg,
			Err:  errors.New(b.Msg),
		}
	}
}

func (b *BaseError) Error() string {
	return b.Err.Error()
}

func (b *BaseError) Unwrap() error {
	return b.Err
}

type ErrExist struct {
	Msg string
	Err error
}

func (m *ErrExist) Error() string {
	return m.Msg
}

type ErrNotFound struct {
	BaseError
}

func NewErrNotFound(msg string, err error) *ErrNotFound {
	e := &ErrNotFound{}
	e.SetCode("default.not_found_error")
	//e.SetMsg(msg)
	e.SetErr(msg, err)

	return e
}

type UnknownError struct {
	BaseError
}

func NewUnknownErr(err error) *ErrNotFound {
	e := &ErrNotFound{}
	e.SetCode("default.unknown_error")
	//e.SetMsg("Unknown error")
	e.SetErr("Unknown error", err)

	return e
}

type InternalErr struct {
	BaseError
}

func NewInternalErr(err error) *InternalErr {
	e := &InternalErr{}
	e.SetCode("default.internal_error")
	e.SetErr("Internal error", err)

	return e
}
