package v1

import (
	"net/http"

	"github.com/Axope/JOJ/common/jwt"
	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/common/response"
	"github.com/Axope/JOJ/internal/service"
	"github.com/gin-gonic/gin"
)

type submitAPI struct {
}

var SubmitAPI = new(submitAPI)

// Submit
//
//	@Tags		Submit
//	@Param		data	body		request.SubmitRequest	true	"提交代码"
//	@Success	200		{object}	response.Response{data=response.SubmitResponse}
//	@Router		/submit [post]
//	@Security	ApiKeyAuth
func (s *submitAPI) Submit(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBindJSON error", log.Any("err", err))
		return
	}
	var err error
	req.UID, err = jwt.GetJWT().GetUserID(c)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("GetUserID error", log.Any("err", err))
		return
	}

	log.LoggerSugar.Infof("Submit:[%+v]", req)

	err = service.SubmitService.Submit(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("change password failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.SubmitResponse{
		Msg: "submit success",
	}))
}
