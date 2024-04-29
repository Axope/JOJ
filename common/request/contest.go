package request

import (
	"time"

	"github.com/Axope/JOJ/internal/model/contest"
)

type GetContestListRequest struct {
	StartIndex int64 `json:"startIndex" form:"startIndex" binging:"required,min=1"`
	Length     int64 `json:"length" form:"length" binging:"required,min=1,max=100"`
}

type GetContestRequest struct {
	CID string `json:"cid" form:"cid"`
}

type CreateContestRequest struct {
	Title     string                   `json:"title" form:"title"`
	Problems  []contest.ContestProblem `json:"problems" form:"problems"`
	StartTime time.Time                `json:"startTime" form:"startTime"`
	Duration  time.Duration            `json:"duration" form:"duration"`
	Note      string                   `json:"note" form:"note"`
	Rule      string                   `json:"rule" form:"rule"`
}
