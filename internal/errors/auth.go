package errors

type InvalidCredentialsError struct {
}

func (e InvalidCredentialsError) Error() string {
	return "invalid credentials"
}

type InvalidJWTError struct {
}

func (e InvalidJWTError) Error() string {
	return "invalid jwt"
}

type InvalidTransmissionError struct {
}

func (e InvalidTransmissionError) Error() string {
	return "invalid transmission type"
}
