package contest

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)



type Standings interface {
	Register(uid uint, problems []ContestProblem) error
	Unregister(uid uint) error
	Accept(uid uint, i int, submitTime, StartTime time.Time) error
	Fail(uid uint, i int) error
	GetStandingsByRank(startIdx int64, len int64) ([][]ProblemSolveStatus, error)
	Close() error
}

func NewStandings(rule string, cid primitive.ObjectID, problemsCnt int) (Standings, error) {
	switch rule {
	case "ACM":
		return NewACMStandings(cid, problemsCnt)
	case "OI":
		// TODO: OI rule
	case "IOI":
		// TODO: IOI rule
	default:
	}
	return nil, fmt.Errorf("unknown contest rule: %v", rule)
}
