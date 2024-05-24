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
	TagsJson        string `json:"tagsJson" form:"tagsJson"`
	InputFormat     string `json:"inputFormat" form:"inputFormat"`
	OutputFormat    string `json:"outputFormat" form:"outputFormat"`
	// TestSamples []model.TestCase `json:"testSamples" form:"testSamples"`
}

// type UpdateProblemRequest struct {
// 	UpdateID
// 	NewProblem model.Problem
// }

// type DeleteProblemRequest struct {
// 	UpdateID
// }
