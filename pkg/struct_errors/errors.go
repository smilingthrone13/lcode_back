package struct_errors

type ErrExist struct {
	Msg string
	Err error
}

func (m *ErrExist) Error() string {
	return m.Msg
}
