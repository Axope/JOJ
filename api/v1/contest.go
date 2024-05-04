package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Axope/JOJ/common/jwt"
	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/common/response"
	"github.com/Axope/JOJ/internal/dao"
	"github.com/Axope/JOJ/internal/manager"
	"github.com/Axope/JOJ/internal/model"
	"github.com/Axope/JOJ/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type contestAPI struct {
}

var ContestAPI = new(contestAPI)

// GetContestList
//
//	@Tags		Contest
//	@Param		data	formData	request.GetContestListRequest	true	"左端点, 长度"
//	@Success	200		{object}	response.Response{data=response.GetContestListResponse}
//	@Router		/contest/getContestList [post]
//	@Security	ApiKeyAuth
func (*contestAPI) GetContestList(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.GetContestListRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBind error", log.Any("err", err))
		return
	}

	log.LoggerSugar.Infof("GetContestList:[%+v]", req)

	uid, err := jwt.GetJWT().GetUserID(c)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("GetUserID error", log.Any("err", err))
		return
	}
	contests, err := service.ContestService.GetContestList(&req, uid)
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

// RegisterContest
//
//	@Tags		Contest
//	@Param		data	body		request.RegisterContestRequest	true	"注册比赛"
//	@Success	200		{object}	response.Response{data=response.RegisterContestResponse}
//	@Router		/contest/register [post]
//	@Security	ApiKeyAuth
func (*contestAPI) RegisterContest(c *gin.Context) {
	defer log.Logger.Sync()
	defer log.LoggerSugar.Sync()

	var req request.RegisterContestRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBindJSON error", log.Any("err", err))
		return
	}
	req.UID, err = jwt.GetJWT().GetUserID(c)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("GetUserID error", log.Any("err", err))
		return
	}

	log.LoggerSugar.Infof("RegisterContest:[cid: %+v]", req)

	err = service.ContestService.RegisterContest(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: RegisterContest failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.RegisterContestResponse{
		Success: true,
		Msg:     "register success",
	}))
}

// UnregisterContest
//
//	@Tags		Contest
//	@Param		data	body		request.UnregisterContestRequest	true	"取消注册比赛"
//	@Success	200		{object}	response.Response{data=response.UnregisterContestResponse}
//	@Router		/contest/unregister [post]
//	@Security	ApiKeyAuth
func (*contestAPI) UnregisterContest(c *gin.Context) {
	defer log.Logger.Sync()
	defer log.LoggerSugar.Sync()

	var req request.UnregisterContestRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBindJSON error", log.Any("err", err))
		return
	}
	req.UID, err = jwt.GetJWT().GetUserID(c)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("GetUserID error", log.Any("err", err))
		return
	}

	log.LoggerSugar.Infof("UnregisterContest:[cid: %+v]", req)

	err = service.ContestService.UnregisterContest(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: UnregisterContest failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.UnregisterContestResponse{
		Success: true,
		Msg:     "unregister success",
	}))
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

	time.AfterFunc(time.Until(req.StartTime), func() {
		log.Logger.Debug("contest start")
		if err := contest.Start(); err != nil {
			panic(err)
		}
	})
	time.AfterFunc(time.Until(req.StartTime.Add(time.Minute*time.Duration(req.Duration))), func() {
		log.Logger.Debug("contest close")
		if err := contest.Close(); err != nil {
			panic(err)
		}
		manager.ContestManager.DelContest(contest.CID.Hex())
	})
}

// GetStandingsByRank
//
//	@Tags		Contest
//	@Param		data	query		request.GetStandingsByRankRequest	true	"获取排名"
//	@Success	200		{object}	response.Response{data=response.GetStandingsByRankResponse}
//	@Router		/contest/getStandingsByRank [get]
func (*contestAPI) GetStandingsByRank(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.GetStandingsByRankRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBindQuery error", log.Any("err", err))
		return
	}
	log.LoggerSugar.Infof("GetStandingsByRank req:[%+v]", req)

	rankList, rule, err := service.ContestService.GetStandingsByRank(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: GetStandingsByRank failed", log.Any("err", err))
		return
	}
	uidList := make([]uint, 0)
	for i := range rankList {
		uidS := rankList[i].Uid
		uid, _ := strconv.ParseUint(uidS, 10, 64)
		uidList = append(uidList, uint(uid))
	}
	var usernameList []string
	var user model.User
	var result *gorm.DB
	for _, uid := range uidList {
		result = dao.GetMysql().Select("username").First(&user, uid)
		if result.Error != nil {
			c.JSON(http.StatusOK, response.FailMsg(result.Error.Error()))
			log.Logger.Warn("uid not find", log.Any("err", err))
			return
		}
		usernameList = append(usernameList, user.Username)
	}
	log.LoggerSugar.Debugf("mysql query uidList(%+v), result(%+v)",
		uidList, usernameList)

	c.JSON(http.StatusOK, response.SuccessMsg(response.GetStandingsByRankResponse{
		RankList:     rankList,
		UsernameList: usernameList,
		Rule:         rule,
	}))
}

// GetContestSubmissionList
//
//	@Tags		Contest
//	@Param		data	query		request.GetContestSubmissionListRequest	true	"cid"
//	@Success	200		{object}	response.Response{data=response.GetContestSubmissionListResponse}
//	@Router		/contest/getContestSubmissionList [get]
func (s *contestAPI) GetContestSubmissionList(c *gin.Context) {
	defer log.Logger.Sync()

	var req request.GetContestSubmissionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("ShouldBindQuery error", log.Any("err", err))
		return
	}

	log.LoggerSugar.Infof("GetContestSubmissionList:[%+v]", req)

	submissions, err := service.ContestService.GetContestSubmissionList(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		log.Logger.Warn("service: GetContestSubmissionList failed", log.Any("err", err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(response.GetContestSubmissionListResponse{
		Submissions: submissions,
	}))
	log.LoggerSugar.Info("service: GetContestSubmissionList:", submissions)
}