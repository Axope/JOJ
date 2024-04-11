package request

import (
	"time"

	"github.com/Axope/JOJ/internal/model"
)

type SubmitRequest struct {
	UID        uint          `json:"-"`
	PID        string        `json:"pid"`
	SubmitTime time.Time     `json:"submitTime"`
	Lang       model.LangSet `json:"lang"`
	SubmitCode string        `json:"submitCode"`
}
