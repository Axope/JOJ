package request

import (
	"github.com/Axope/JOJ/internal/model"
)

type JudgeRequest struct {
	SID         string           `json:"sid"`
	TimeLimit   int64            `json:"timeLimit"`
	MemoryLimit int64            `json:"memoryLimit"`
	TestSamples []model.TestCase `json:"testSamples"`
	TestCases   []model.TestCase `json:"testCases"`
	Lang        model.LangSet    `json:"lang"`
	SubmitCode  string           `json:"submitCode"`
}
