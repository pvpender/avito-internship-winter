package errors

type ShutdownError struct {
}

func (err ShutdownError) Error() string {
	return "ShutdownError: server is shutting down"
}
