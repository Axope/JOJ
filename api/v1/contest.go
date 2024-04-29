package v1

import (
	"net/http"
	"time"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/common/response"
	"github.com/Axope/JOJ/internal/service"
	"github.com/gin-gonic/gin"
)

type contestAPI struct {
}

var ContestAPI = new(contestAPI)

// GetContestList
//
//	@Tags		Contest
//	@Param		data	query		request.GetContestListRequest	true	"左端点, 长度"
//	@Success	200		{object}	response.Response{data=response.GetContestListResponse}
//	@Router		/contest/getContestList [get]
func (*contestAPI) GetContestList(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.GetContestListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBindQuery error", log.Any("err", err))
		return
	}

	log.LoggerSugar.Infof("GetContestList:[%+v]", req)

	contests, err := service.ContestService.GetContestList(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: GetContestList failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.GetContestListResponse{
		Contests: contests,
	}))
	log.LoggerSugar.Info("service: GetContestList:", contests)
}

// GetContest
//
//	@Tags		Contest
//	@Param		cid	path		string	true	"contest ID"
//	@Success	200	{object}	response.Response{data=response.GetContestResponse}
//	@Router		/contest/{cid} [get]
func (*contestAPI) GetContest(c *gin.Context) {
	defer log.Logger.Sync()

	cid := c.Param("cid")
	if cid == "" {
		c.JSON(http.StatusOK, response.FailMsg("cid empty"))
		log.Logger.Warn("cid empty")
		return
	}

	log.LoggerSugar.Infof("GetContest:[cid: %+v]", cid)

	contest, err := service.ContestService.GetContest(cid)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: GetContest failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.GetContestResponse{
		Contest: *contest,
	}))
	log.LoggerSugar.Info("service: GetContest:", contest)
}

// CreateContest
//
//	@Accept		multipart/form-data
//	@Tags		Contest
//	@Param		data	formData	request.CreateContestRequest	true	"contest"
//	@Success	200		{object}	response.Response{data=response.CreateContestResponse}
//	@Router		/contest/createContest [post]
//	@Security	ApiKeyAuth
func (*contestAPI) CreateContest(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.CreateContestRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBind error", log.Any("err", err))
		return
	}
	log.LoggerSugar.Infof("CreateContest req:[%+v]", req)

	contest, err := service.ContestService.CreateContest(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: CreateContest failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.CreateContestResponse{
		Success: true,
		Msg:     "create contest success",
	}))

	curTime := time.Now()
	startDura := req.StartTime.Sub(curTime)
	stopDura := req.StartTime.Add(req.Duration).Sub(curTime)
	go func() {
		<-time.After(startDura)
		contest.Start()
	}()
	go func() {
		<-time.After(stopDura)
		contest.Close()
	}()

}

