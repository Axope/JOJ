package v1

import (
	"net/http"

	"github.com/Axope/JOJ/common/jwt"
	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/common/response"
	"github.com/Axope/JOJ/internal/model"
	"github.com/Axope/JOJ/internal/service"
	"github.com/gin-gonic/gin"
)

type userAPI struct {
}

var UserAPI = new(userAPI)

// Register
//
//	@Tags		User
//	@Param		data	formData	request.RegisterRequest	true	"用户名, 密码"
//	@Success	200		{object}	response.Response{data=response.RegisterResponse}
//	@Router		/user/register [post]
func (u *userAPI) Register(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBind error", log.Any("err", err))
		return
	}

	log.LoggerSugar.Infof("register:[%v]", req)

	user, err := service.UserService.Register(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("register failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.RegisterResponse{
		Username: user.Username,
	}))
	log.LoggerSugar.Infof("user(%v) registration success", user)
}

func createToken(user *model.User) (string, error) {
	return jwt.GetJWT().CreateToken(*request.NewCustomClaims(user.UUID, user.ID, user.Admin, jwt.GetJWT().Expire))
}

// Login
//
//	@Tags		User
//	@Param		data	formData	request.LoginRequest	true	"用户名, 密码"
//	@Success	200		{object}	response.Response{data=response.LoginResponse}
//	@Router		/user/login [post]
func (u *userAPI) Login(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBind error", log.Any("err", err))
		return
	}

	log.LoggerSugar.Infof("login:[%v]", req)

	user, err := service.UserService.Login(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("login failed", log.Any("err", err))
		return
	}

	token, err := createToken(user)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("create token failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.LoginResponse{
		UID:      user.ID,
		Username: user.Username,
		Token:    token,
	}))
	log.LoggerSugar.Infof("user(%v) login success", user)
}

// ChangePassword
//
//	@Tags		User
//	@Param		data	formData	request.ChangePasswordRequest	true	"用户名, 密码, 新密码"
//	@Success	200		{object}	response.Response{data=response.ChangePasswordResponse}
//	@Router		/user/changePassword [post]
//	@Security	ApiKeyAuth
func (u *userAPI) ChangePassword(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.ChangePasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBind error", log.Any("err", err))
		return
	}
	var err error
	req.ID, err = jwt.GetJWT().GetUserID(c)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("GetUserID error", log.Any("err", err))
		return
	}

	log.LoggerSugar.Infof("change password:[%v]", req)

	err = service.UserService.ChangePassword(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("change password failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.ChangePasswordResponse{
		Msg: "modify success",
	}))
	log.LoggerSugar.Infof("user(%v) change password success", req)
}
