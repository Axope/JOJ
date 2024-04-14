package request

type GetProblemListRequest struct {
	StartIndex int64 `json:"startIndex" form:"startIndex" binging:"required,min=1"`
	Length     int64 `json:"length" form:"length" binging:"required,min=1,max=100"`
}

type GetProblemRequest struct {
	PID string `json:"pid" form:"pid"`
}

type CreateProblemRequest struct {
	Title           string `json:"title" form:"title"`
	TimeLimit       int64  `json:"timeLimit" form:"timeLimit"`
	MemoryLimit     int64  `json:"memoryLimit" form:"memoryLimit"`
	Description     string `json:"description" form:"description"`
	TestSamplesJson string `json:"testSamplesJson" form:"testSamplesJson"`
}

// type UpdateProblemRequest struct {
// 	UpdateID
// 	NewProblem model.Problem
// }

// type DeleteProblemRequest struct {
// 	UpdateID
// }
