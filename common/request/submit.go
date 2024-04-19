package request

import (
	"time"
)

type SubmitRequest struct {
	UID        uint      `json:"-"`
	PID        string    `json:"pid"`
	SubmitTime time.Time `json:"submitTime"`
	Lang       int32     `json:"lang"`
	SubmitCode string    `json:"submitCode"`
}
