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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContestProblem struct {
	PID   primitive.ObjectID `bson:"pid" json:"pid"`
	Nick  string             `bson:"nick" json:"nick"`
	Title string             `bson:"title" json:"title"`
}
type Contest struct {
	CID        primitive.ObjectID `bson:"_id,omitempty" json:"cid"`
	Title      string             `bson:"title" json:"title"`
	Status     ContestStatus      `bson:"status" json:"status"`
	Standings  Standings          `bson:"-" json:"-"`
	Problems   []ContestProblem   `bson:"problems" json:"problems"`
	StartTime  time.Time          `bson:"startTime" json:"startTime"`
	Duration   time.Duration      `bson:"duration" json:"duration"`
	Note       string             `bson:"note" json:"note"`
	Registered []uint             `bson:"registered,omitempty" json:"registered,omitempty"`
	rwLock     *sync.RWMutex      `bson:"-" json:"-"`
	Rule       string             `bson:"rule" json:"rule"`

	RankList []RankListData `bson:"rankList,omitempty" json:"-"` // 不直接存储
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
		Rule:      rule,
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

	return c, nil
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
func (c *Contest) GetStandingsByRank(startIdx int64, len int64) ([]RankListData, error) {
	return c.Standings.GetStandingsByRank(startIdx, len)
}
func (c *Contest) Start() error {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	c.Status = RUNNING

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: c.Status}}}}
	if _, err := dao.GetContestColl().UpdateByID(context.TODO(), c.CID, update); err != nil {
		return err
	}

	findOptions := options.FindOne()
	findOptions.SetProjection(bson.D{{Key: "registered", Value: 1}})
	var contest Contest
	err := dao.GetContestColl().FindOne(context.TODO(), bson.D{{Key: "_id", Value: c.CID}}, findOptions).Decode(&contest)
	if err != nil {
		log.Logger.Error("FindOne error", log.Any("_id", c.CID), log.Any("findOptions", findOptions))
		return err
	}

	for _, uid := range contest.Registered {
		if err := c.Standings.Register(uid, c.Problems); err != nil {
			return err
		}
	}
	return nil
}
func (c *Contest) Close() error {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	c.Status = CLOSE

	if err := c.Standings.Close(); err != nil {
		log.Logger.Error("close error", log.Any("err", err))
		return err
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: c.Status}}}}
	if _, err := dao.GetContestColl().UpdateByID(context.TODO(), c.CID, update); err != nil {
		return err
	}
	return nil
}
