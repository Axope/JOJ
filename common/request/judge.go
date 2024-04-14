package request

import (
	"github.com/Axope/JOJ/internal/model"
)

type JudgeRequest struct {
	SID         string        `json:"sid"`
	PID         string        `json:"pid"`
	TimeLimit   int64         `json:"timeLimit"`
	MemoryLimit int64         `json:"memoryLimit"`
	TestCases   []string      `json:"testCases"`
	Lang        model.LangSet `json:"lang"`
	SubmitCode  string        `json:"submitCode"`
}
