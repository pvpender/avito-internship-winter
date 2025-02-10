package errors

import "fmt"

type LoadEnvError struct {
	EnvName string
}

func (e *LoadEnvError) Error() string {
	return fmt.Sprintf("LoadEnvError: environment variable %s is not set", e.EnvName)
}
