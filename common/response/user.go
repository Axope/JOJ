package response

type RegisterResponse struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	UID      uint   `json:"uid"`
	Username string `json:"username"`
	Admin    int    `json:"admin"`
	Token    string `json:"token"`
}

type ChangePasswordResponse struct {
	Msg string `json:"msg"`
}
