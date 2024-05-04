package contest

import (
	"github.com/Axope/JOJ/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProblemSolveStatus struct {
	PID     primitive.ObjectID `bson:"_id" json:"pid"`
	Nick    string             `bson:"nick" json:"nick"`
	Status  model.StatusSet    `bson:"status" json:"status"`
	Penalty int64              `bson:"penalty" json:"penalty"`
	FailCnt int64              `bson:"failCnt" json:"failCnt"`
}

func NewProblemSolveStatus(problems []ContestProblem) []ProblemSolveStatus {
	s := make([]ProblemSolveStatus, len(problems))
	for i, p := range problems {
		s[i].PID = p.PID
		s[i].Nick = p.Nick
		s[i].Status = model.UNSUBMIT
		s[i].Penalty = 0
		s[i].FailCnt = 0
	}
	return s
}
func calc(problemSolveStatus []ProblemSolveStatus) (int, int64) {
	var ACCnt int
	var totPenalty int64
	for _, ps := range problemSolveStatus {
		if ps.Status == model.AC {
			ACCnt++
			totPenalty += ps.Penalty
			totPenalty += ps.FailCnt * 20
		}
	}
	return ACCnt, totPenalty
}
