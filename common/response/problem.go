package response

import "github.com/Axope/JOJ/internal/model"

type SimpleProblem struct {
	PID         string   `json:"pid" bson:"_id"`
	Title       string   `json:"title"`
	TimeLimit   int64    `json:"timeLimit"`
	MemoryLimit int64    `json:"memoryLimit"`
	Tags        []string `json:"tags"`
}
type GetProblemListResponse struct {
	Problems []SimpleProblem `json:"problems"`
	Total    int64           `json:"total"`
}

type GetProblemResponse struct {
	Problem model.Problem `json:"problem"`
}

type CreateProblemResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"message"`
}

type UploadDatasResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"message"`
}

type GetTagsResponse struct {
	Tags []string `json:"tags"`
}
