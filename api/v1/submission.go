package v1

import (
	"net/http"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/common/response"
	"github.com/Axope/JOJ/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type submissionAPI struct {
}

var Submission = new(submissionAPI)

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
		log.Logger.Warn("ShouldBindQuery error", zap.Any("err", err))
		return
	}

	log.LoggerSuger.Infof("GetSubmissionList:[%+v]", req)

	submissions, err := service.SubmissionService.GetSubmissionList(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: GetSubmissionList failed", zap.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.GetSubmissionListResponse{
		Submissions: submissions,
	}))
	log.LoggerSuger.Info("service: GetSubmissionList:", submissions)
}
