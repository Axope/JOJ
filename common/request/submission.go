package request

type GetSubmissionListRequest struct {
	UID uint   `json:"uid" form:"uid"`
	PID string `json:"pid" form:"pid"`
}
