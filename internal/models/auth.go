package models

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (req *AuthRequest) Validate() bool {
	if len(req.Username) == 0 || len(req.Password) == 0 {
		return false
	}

	return true
}
