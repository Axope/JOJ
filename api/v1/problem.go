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

type problemAPI struct {
}

var ProblemAPI = new(problemAPI)

// GetProblemList
//
//	@Tags		Problem
//	@Param		data	query		request.GetProblemListRequest	true	"左端点, 长度"
//	@Success	200		{object}	response.Response{data=response.GetProblemListResponse}
//	@Router		/problem/getProblemList [get]
func (u *problemAPI) GetProblemList(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.GetProblemListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBindQuery error", zap.Any("err", err))
		return
	}

	log.LoggerSuger.Infof("GetProblemList:[%+v]", req)

	problems, err := service.ProblemService.GetProblemList(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: GetProblemList failed", zap.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.GetProblemListResponse{
		Problems: problems,
	}))
	log.LoggerSuger.Info("service: GetProblemList:", problems)
}

// GetProblem
//
//	@Tags		Problem
//	@Param		data	query		request.GetProblemRequest	true	"problem ID"
//	@Success	200		{object}	response.Response{data=response.GetProblemResponse}
//	@Router		/problem/getProblem [get]
func (u *problemAPI) GetProblem(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.GetProblemRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBindQuery error", zap.Any("err", err))
		return
	}

	log.LoggerSuger.Infof("GetProblem:[%+v]", req)

	problem, err := service.ProblemService.GetProblem(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: GetProblem failed", zap.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.GetProblemResponse{
		Problem: *problem,
	}))
	log.LoggerSuger.Info("service: GetProblem:", problem)
}

// // CreateProblem
// //
// //	@Tags		Problem
// //	@Param		data	formData	request.CreateProblemRequest	true	"problem"
// //	@Success	200		{object}	response.Response{data=response.CreateProblemResponse}
// //	@Router		/problem/createProblem [post]
// func (u *problemAPI) CreateProblem(c *gin.Context) {

// }

// // UpdateProblem
// //
// //	@Tags		Problem
// //	@Param		data	formData		request.UpdateProblemRequest	true	"problem id, problem"
// //	@Success	200		{object}	response.Response{data=response.UpdateProblemResponse}
// //	@Router		/problem/updateProblem [put]
// func (u *problemAPI) UpdateProblem(c *gin.Context) {

// }

// // DeleteProblem
// //
// //	@Tags		Problem
// // //	@Param		data	query		request.DeleteProblemRequest	true	"problem id"
// //	@Success	200		{object}	response.Response{data=response.DeleteProblemResponse}
// //	@Router		/problem/deleteProblem [delete]
// func (u *problemAPI) DeleteProblem(c *gin.Context) {

// }
