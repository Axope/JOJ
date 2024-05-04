package request

import (
	"time"
)

type SubmitRequest struct {
	UID        uint      `json:"-"`
	PID        string    `json:"pid"`
	CID        string    `json:"cid"`
	SubmitTime time.Time `json:"submitTime"`
	Lang       int32     `json:"lang"`
	SubmitCode string    `json:"submitCode"`
}
