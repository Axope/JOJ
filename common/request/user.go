package request

type RegisterRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type ChangePasswordRequest struct {
	ID          uint   `json:"-"`
	Password    string `json:"password" form:"password"`
	NewPassword string `json:"newPassword" form:"newPassword"`
}
