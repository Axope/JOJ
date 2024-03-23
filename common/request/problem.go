package request

import "github.com/Axope/JOJ/internal/model"

// "github.com/Axope/JOJ/internal/model"
// "go.mongodb.org/mongo-driver/bson/primitive"

type GetProblemListRequest struct {
	StartIndex int64 `json:"startIndex" form:"startIndex" binging:"required,min=1"`
	Length     int64 `json:"length" form:"length" binging:"required,min=1,max=100"`
}

type GetProblemRequest struct {
	PID string `json:"pid" form:"pid"`
}

type CreateProblemRequest struct {
	Problem model.Problem `json:"problem"`
}

// type UpdateProblemRequest struct {
// 	UpdateID
// 	NewProblem model.Problem
// }

// type DeleteProblemRequest struct {
// 	UpdateID
// }
