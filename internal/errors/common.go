package errors

type NilPointerError struct {
}

func (e NilPointerError) Error() string {
	return "pointer is nil"
}

type ValidationError struct {
}

func (e ValidationError) Error() string {
	return "validation failed"
}
