package errors

type InvalidCredentials struct {
}

func (e InvalidCredentials) Error() string {
	return "invalid credentials"
}

type InvalidJWT struct {
}

func (e InvalidJWT) Error() string {
	return "invalid jwt"
}
