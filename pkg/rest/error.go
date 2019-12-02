package rest

type Error struct {
	Msg string
	Err error
}

func (e *Error) Error() string {
	if e.Err != nil {
		if e.Msg != "" {
			return e.Msg + ": " + e.Err.Error()
		}
		return "rest: " + e.Err.Error()
	}
	return e.Msg
}
