package response

import "github.com/Axope/JOJ/internal/model"

type GetProblemListResponse struct {
	Problems []model.Problem `json:"problems"`
}

type GetProblemResponse struct {
	Problem model.Problem `json:"problem"`
}
