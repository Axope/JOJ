package v1

import (
	"net/http"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/common/response"
	"github.com/Axope/JOJ/configs"
	"github.com/Axope/JOJ/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		log.Logger.Warn("ShouldBindQuery error", log.Any("err", err))
		return
	}

	log.LoggerSugar.Infof("GetProblemList:[%+v]", req)

	total, problems, err := service.ProblemService.GetProblemList(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: GetProblemList failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.GetProblemListResponse{
		Problems: problems,
		Total: total,
	}))
	log.LoggerSugar.Info("service: GetProblemList:", problems)
}

// GetProblem
//
//	@Tags		Problem
//	@Param		pid	path		string	true	"problem ID"
//	@Success	200	{object}	response.Response{data=response.GetProblemResponse}
//	@Router		/problem/{pid} [get]
func (u *problemAPI) GetProblem(c *gin.Context) {
	defer log.Logger.Sync()

	pid := c.Param("pid")
	if pid == "" {
		c.JSON(http.StatusOK, response.FailMsg("pid empty"))
		log.Logger.Warn("pid empty")
		return
	}

	log.LoggerSugar.Infof("GetProblem:[pid: %+v]", pid)

	problem, err := service.ProblemService.GetProblem(pid)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: GetProblem failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.GetProblemResponse{
		Problem: *problem,
	}))
	log.LoggerSugar.Info("service: GetProblem:", problem)
}

// GetTags
//
//	@Tags		Problem
//	@Success	200	{object}	response.Response{data=response.GetTagsResponse}
//	@Router		/problem/getTags [get]
func (u *problemAPI) GetTags(c *gin.Context) {
	c.JSON(http.StatusOK, response.SuccessMsg(response.GetTagsResponse{
		Tags: configs.GetTagsList(),
	}))
}

// CreateProblem
//
//	@Accept		multipart/form-data
//	@Tags		Problem
//	@Param		data	formData	request.CreateProblemRequest	true	"problem"
//	@Param		file	formData	file							true	"文件"
//	@Success	200		{object}	response.Response{data=response.CreateProblemResponse}
//	@Router		/problem/createProblem [post]
//	@Security	ApiKeyAuth
func (u *problemAPI) CreateProblem(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.CreateProblemRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBind error", log.Any("err", err))
		return
	}
	log.LoggerSugar.Infof("CreateProblem req:[%+v]", req)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("file parse error", log.Any("err", err))
		return
	}
	// save zip
	pid := primitive.NewObjectID()
	zipPath := configs.GetDatasConfig().DirPath + pid.Hex() + ".zip"
	if err := c.SaveUploadedFile(file, zipPath); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailMsg(err.Error()))
		log.Logger.Warn("SaveUploadedFile error", log.Any("err", err))
		return
	}
	log.Logger.Info("CreateProblem save file ok")

	testCasesPath := configs.GetDatasConfig().DirPath + pid.Hex() + "/"
	if err := service.ProblemService.CreateProblem(&req, pid, zipPath, testCasesPath); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: CreateProblem failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.CreateProblemResponse{
		Success: true,
		Msg:     "create problem success",
	}))
}

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
