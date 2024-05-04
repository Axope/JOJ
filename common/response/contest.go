package response

import (
	"time"

	"github.com/Axope/JOJ/internal/model"
	"github.com/Axope/JOJ/internal/model/contest"
)

type SimpleContest struct {
	CID          string                `json:"cid" bson:"_id"`
	Title        string                `json:"title" bson:"title"`
	Status       contest.ContestStatus `json:"status" bson:"status"`
	StartTime    time.Time             `json:"startTime" bson:"startTime"`
	Duration     time.Duration         `json:"duration" bson:"duration"`
	Note         string                `json:"note" bson:"note"`
	Registered   []uint                `bson:"registered" json:"-"`
	IsRegistered bool                  `json:"isRegistered" bson:"isRegistered"`
}
type GetContestListResponse struct {
	Contests []SimpleContest `json:"contests"`
}

type GetContestResponse struct {
	Contest contest.Contest `json:"contest"`
}

type CreateContestResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"message"`
}

type RegisterContestResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"message"`
}

type UnregisterContestResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"message"`
}

type GetStandingsByRankResponse struct {
	RankList     []contest.RankListData `json:"rankList"`
	UsernameList []string               `json:"usernameList"`
	Rule         string                 `json:"rule"`
}

type GetContestSubmissionListResponse struct {
	Submissions []model.Submission `json:"submissions"`
}
