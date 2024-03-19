package v1

import (
	"github.com/Axope/JOJ/common/jwt"
	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/common/response"
	"github.com/Axope/JOJ/internal/model"
	"github.com/Axope/JOJ/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type userAPI struct {
}

var UserAPI = new(userAPI)

// Register
//
//	@Tags		User
//	@Param		data	body		request.RegisterRequest	true	"用户名, 密码"
//	@Success	200		{object}	response.Response{data=response.RegisterResponse}
//	@Router		/user/register [post]
func (u *userAPI) Register(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBindJSON error", zap.Any("err", err))
		return
	}

	log.LoggerSuger.Infof("register:[%v]", req)

	user, err := service.UserService.Register(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("register failed", zap.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.RegisterResponse{
		Username: user.Username,
	}))
	log.LoggerSuger.Infof("user(%v) registration success", user)
}

func createToken(user *model.User) (string, error) {
	return jwt.GetJWT().CreateToken(*request.NewCustomClaims(user.UUID, user.ID, 0, jwt.GetJWT().Expire))
}

// Login
//
//	@Tags		User
//	@Param		data	body		request.LoginRequest	true	"用户名, 密码"
//	@Success	200		{object}	response.Response{data=response.LoginResponse}
//	@Router		/user/login [post]
func (u *userAPI) Login(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBindJSON error", zap.Any("err", err))
		return
	}

	log.LoggerSuger.Infof("login:[%v]", req)

	user, err := service.UserService.Login(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("login failed", zap.Any("err", err))
		return
	}

	token, err := createToken(user)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("create token failed", zap.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.LoginResponse{
		Username: user.Username,
		Token:    token,
	}))
	log.LoggerSuger.Infof("user(%v) login success", user)
}

// ChangePassword
//
//	@Tags		User
//	@Param		data	body		request.ChangePasswordRequest	true	"用户名, 密码, 新密码"
//	@Success	200		{object}	response.Response{data=response.ChangePasswordResponse}
//	@Router		/user/changePassword [post]
//	@Security	ApiKeyAuth
func (u *userAPI) ChangePassword(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBindJSON error", zap.Any("err", err))
		return
	}
	var err error
	req.ID, err = jwt.GetJWT().GetUserID(c)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("GetUserID error", zap.Any("err", err))
		return
	}

	log.LoggerSuger.Infof("change password:[%v]", req)

	err = service.UserService.ChangePassword(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("change password failed", zap.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.ChangePasswordResponse{
		Msg: "modify success",
	}))
	log.LoggerSuger.Infof("user(%v) change password success", req)
}
