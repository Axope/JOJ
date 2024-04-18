package response

import "github.com/Axope/JOJ/internal/model"

type SimpleProblem struct {
	PID   string `json:"pid" bson:"_id"`
	Title string `json:"title"`
}
type GetProblemListResponse struct {
	Problems []SimpleProblem `json:"problems"`
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
