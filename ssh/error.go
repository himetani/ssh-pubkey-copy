package ssh

// SkipCopyError is struct of SkipCopyError
type SkipCopyError struct {
	msg string
}

func (s *SkipCopyError) Error() string {
	return s.msg
}
