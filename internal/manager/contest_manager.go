package manager

import (
	"errors"
	"sync"

	"github.com/Axope/JOJ/internal/model/contest"
)

type contestManager struct {
	contestMap sync.Map // key:cid, value:*contest.Contest
}

var ContestManager contestManager

func (cm *contestManager) NewContest(c *contest.Contest) {
	cm.contestMap.Store(c.CID.Hex(), c)
}
func (cm *contestManager) DelContest(cid string) {
	cm.contestMap.Delete(cid)
}
func (cm *contestManager) GetContest(cid string) (*contest.Contest, error) {
	v, ok := cm.contestMap.Load(cid)
	if !ok {
		return nil, errors.New("map load error")
	}
	c, ok := v.(*contest.Contest)
	if !ok {
		return nil, errors.New("type assert error")
	}
	return c, nil 
}
