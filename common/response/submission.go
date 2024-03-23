package response

import "github.com/Axope/JOJ/internal/model"

type GetSubmissionListResponse struct {
	Submissions []model.Submission `json:"submissions"`
}