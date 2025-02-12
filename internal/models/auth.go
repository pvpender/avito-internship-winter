package models

import "errors"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func NewAuthRequest(username string, password string) (*AuthRequest, error) {
	if len(username) == 0 || len(password) == 0 {
		return nil, errors.New("username and password cannot be empty")
	}

	return &AuthRequest{Username: username, Password: password}, nil
}

func NewAuthResponse(token string) (*AuthResponse, error) {
	if len(token) == 0 {
		return nil, errors.New("token cannot be empty")
	}

	return &AuthResponse{Token: token}, nil
}
