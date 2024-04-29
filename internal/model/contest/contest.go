package contest

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/internal/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContestProblem struct {
	PID   primitive.ObjectID `bson:"pid" json:"pid"`
	Nick  string             `bson:"nick" json:"nick"`
	Title string             `bson:"title" json:"title"`
}
type Contest struct {
	CID       primitive.ObjectID `bson:"_id,omitempty" json:"cid"`
	Title     string             `bson:"title" json:"title"`
	Status    ContestStatus      `bson:"status" json:"status"`
	Standings Standings          `bson:"-" json:"-"`
	Problems  []ContestProblem   `bson:"problems" json:"problems"`
	StartTime time.Time          `bson:"startTime" json:"startTime"`
	Duration  time.Duration      `bson:"duration" json:"duration"`
	Note      string             `bson:"note" json:"note"`
	// RankList  RankList           `bson:"rankList" json:"rankList"`
	rwLock *sync.RWMutex `bson:"-" json:"-"`
}

func NewContest(cid primitive.ObjectID, title string, problems []ContestProblem,
	startTime time.Time, duration time.Duration, note string, rule string) (*Contest, error) {
	c := &Contest{
		CID:       cid,
		Title:     title,
		Problems:  problems,
		StartTime: startTime,
		Duration:  duration,
		Note:      note,
		rwLock:    &sync.RWMutex{},
	}
	standings, err := NewStandings(rule, c.CID, len(c.Problems))
	if err != nil {
		return nil, err
	}
	c.Standings = standings

	now := time.Now()
	if now.Before(startTime) {
		c.Status = REGISTER
	} else if now.Before(startTime.Add(duration)) {
		c.Status = RUNNING
	} else {
		c.Status = CLOSE
	}

	// TODO: rank list store to stable
	return c, nil
}
func (c *Contest) Register(uid uint) error {
	defer log.Logger.Sync()
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	if c.Status == RUNNING {
		log.LoggerSugar.Debugf("uid(%v) register failed, contest is running", uid)
		return fmt.Errorf("contest is running")
	}
	if c.Status == CLOSE {
		log.LoggerSugar.Debugf("uid(%v) register failed, contest is close", uid)
		return fmt.Errorf("contest is close")
	}

	log.LoggerSugar.Debugf("uid(%v) register success", uid)
	c.Standings.Register(uid, c.Problems)

	return nil
}
func (c *Contest) Unregister(uid uint) error {
	defer log.Logger.Sync()
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	if c.Status == RUNNING {
		log.LoggerSugar.Debugf("uid(%v) unregister failed, contest is running", uid)
		return fmt.Errorf("contest is running")
	}
	if c.Status == CLOSE {
		log.LoggerSugar.Debugf("uid(%v) unregister failed, contest is close", uid)
		return fmt.Errorf("contest is close")
	}

	c.Standings.Unregister(uid)
	return nil
}
func (c *Contest) Accept(uid uint, pid primitive.ObjectID, submitTime time.Time) error {
	defer log.Logger.Sync()
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	if c.Status == CLOSE {
		return nil
	}
	if c.Status == REGISTER {
		log.LoggerSugar.Errorf("ERROR: contest is register", uid)
		return fmt.Errorf("contest is register")
	}

	for i, p := range c.Problems {
		if p.PID == pid {
			return c.Standings.Accept(uid, i, submitTime, c.StartTime)
		}
	}
	return fmt.Errorf("pid(%v) is not part of contest(%v: %v)", pid, c.CID, c.Title)
}
func (c *Contest) Fail(uid uint, pid primitive.ObjectID) error {
	defer log.Logger.Sync()
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	if c.Status == CLOSE {
		return nil
	}
	if c.Status == REGISTER {
		log.LoggerSugar.Errorf("ERROR: contest is register", uid)
		return fmt.Errorf("contest is register")
	}

	for i, p := range c.Problems {
		if p.PID == pid {
			return c.Standings.Fail(uid, i)
		}
	}
	return fmt.Errorf("pid(%v) is not part of contest(%v: %v)", pid, c.CID, c.Title)
}
func (c *Contest) GetStandingsByRank(startIdx int64, len int64) ([][]ProblemSolveStatus, error) {
	return c.Standings.GetStandingsByRank(startIdx, len)
}
func (c *Contest) Start() {
	// TODO: review
	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	c.Status = RUNNING
	go func() {
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: c.Status}}}}
		dao.GetContestColl().UpdateByID(context.TODO(), c.CID, update)
	}()
}
func (c *Contest) Close() {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	
	c.Status = CLOSE
	go func() {
		c.Standings.Close()
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: c.Status}}}}
		dao.GetContestColl().UpdateByID(context.TODO(), c.CID, update)
	}()
}
