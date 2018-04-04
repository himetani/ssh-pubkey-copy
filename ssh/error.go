package ssh

type ConnectionError struct {
	msg string
}

func (e *ConnectionError) Error() string {
	return e.msg
}
