package v1

import (
	"net/http"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/common/response"
	"github.com/Axope/JOJ/internal/service"
	"github.com/gin-gonic/gin"
)

type submissionAPI struct {
}

var SubmissionAPI = new(submissionAPI)

// GetSubmissionList
//
//	@Tags		Submission
//	@Param		data	query		request.GetSubmissionListRequest	true	"uid, pid"
//	@Success	200		{object}	response.Response{data=response.GetSubmissionListResponse}
//	@Router		/submission/getSubmissionList [get]
func (s *submissionAPI) GetSubmissionList(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.GetSubmissionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBindQuery error", log.Any("err", err))
		return
	}

	log.LoggerSugar.Infof("GetSubmissionList:[%+v]", req)

	submissions, err := service.SubmissionService.GetSubmissionList(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: GetSubmissionList failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.GetSubmissionListResponse{
		Submissions: submissions,
	}))
	log.LoggerSugar.Info("service: GetSubmissionList:", submissions)
}
