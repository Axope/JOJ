package request

type RegisterRequest struct {
	Username string `json:"username" example:"用户名"`
	Password string `json:"password" example:"密码"`
}

type LoginRequest struct {
	Username string `json:"username" example:"用户名"`
	Password string `json:"password" example:"密码"`
}

type ChangePasswordRequest struct {
	ID          uint   `json:"-"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}
