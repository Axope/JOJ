package response

import (
	"time"

	"github.com/Axope/JOJ/internal/model/contest"
)

type SimpleContest struct {
	CID       string                `json:"cid" bson:"_id"`
	Title     string                `json:"title"`
	Status    contest.ContestStatus `json:"status"`
	StartTime time.Time             `json:"startTime"`
	Duration  time.Duration         `json:"duration"`
	// options
	Note string `json:"note"`
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
