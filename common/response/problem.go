package response

import "github.com/Axope/JOJ/internal/model"

type GetProblemListResponse struct {
	Problems []model.Problem `json:"problems"`
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
