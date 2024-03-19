package response

type RegisterResponse struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type ChangePasswordResponse struct {
	Msg string `json:"msg"`
}
