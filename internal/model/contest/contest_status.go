package contest

type ContestStatus string

const (
	REGISTER ContestStatus = "Register"
	RUNNING  ContestStatus = "Running"
	CLOSE    ContestStatus = "Close"
)
