package request

import (
	"time"
)

type GetContestListRequest struct {
	StartIndex int64 `json:"startIndex" form:"startIndex" binging:"required,min=1"`
	Length     int64 `json:"length" form:"length" binging:"required,min=1,max=100"`
}

type GetContestRequest struct {
	CID string `json:"cid" form:"cid"`
}

type CreateContestRequest struct {
	Title        string    `json:"title" form:"title"`
	ProblemsJson string    `json:"problemsJson" form:"problemsJson"`
	StartTime    time.Time `json:"startTime" form:"startTime"`
	Duration     int64     `json:"duration" form:"duration"`
	Note         string    `json:"note" form:"note"`
	Rule         string    `json:"rule" form:"rule"`
}

type RegisterContestRequest struct {
	CID string `json:"cid"`
	UID uint   `json:"-"`
}

type UnregisterContestRequest struct {
	CID string `json:"cid"`
	UID uint   `json:"-"`
}

type GetStandingsByRankRequest struct {
	Cid      string `json:"cid" form:"cid"`
	StartIdx int64  `json:"startIdx" form:"startIdx"`
	Len      int64  `json:"len" form:"len"`
}

type GetContestSubmissionListRequest struct {
	CID string `json:"cid" form:"cid"`
}
