package contest

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/internal/dao"
	"github.com/Axope/JOJ/internal/middleware/redis"
	"github.com/Axope/JOJ/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const maxRetries = 100

type ACMStandings struct {
	BindCID        primitive.ObjectID
	ProblemsCnt    int
	InfoPrefix     string // $cid:$uid
	SolveSetPrefix string // $cid:$solveCnt
	ZKeys          []string
}

func NewACMStandings(cid primitive.ObjectID, problemsCnt int) (*ACMStandings, error) {
	acm := &ACMStandings{
		BindCID:        cid,
		ProblemsCnt:    problemsCnt,
		InfoPrefix:     cid.Hex() + ":",
		SolveSetPrefix: cid.Hex() + ":",
	}
	for i := 0; i <= problemsCnt; i++ {
		acm.ZKeys = append(acm.ZKeys, acm.SolveSetPrefix+strconv.Itoa(i))
	}
	return acm, nil
}

func (acm *ACMStandings) Register(uid uint, problems []ContestProblem) error {
	defer log.Logger.Sync()
	defer log.LoggerSugar.Sync()
	uidS := strconv.FormatUint(uint64(uid), 10)
	pssJson, err := json.Marshal(NewProblemSolveStatus(problems))
	if err != nil {
		log.Logger.Error("json Marshal error", log.Any("err", err))
		return err
	}
	log.Logger.Debug("get json", log.Any("json", string(pssJson)))
	// update
	ZeroSolveSet := acm.SolveSetPrefix + "0"
	InfoKey := acm.InfoPrefix + uidS

	ctx := context.TODO()
	txf := func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if err := pipe.Set(ctx, InfoKey, pssJson, 0).Err(); err != nil {
				return err
			}
			if err := pipe.ZAdd(ctx, ZeroSolveSet, redis.Z{Score: 0, Member: uidS}).Err(); err != nil {
				return err
			}
			return nil
		})
		return err
	}
	return doTransactional(ctx, txf, InfoKey)
}
func (acm *ACMStandings) Unregister(uid uint) error {
	uidS := strconv.FormatUint(uint64(uid), 10)
	ZeroSolveSet := acm.SolveSetPrefix + "0"
	InfoKey := acm.InfoPrefix + uidS

	ctx := context.TODO()
	txf := func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if err := pipe.Del(ctx, InfoKey).Err(); err != nil {
				return err
			}
			if err := pipe.ZRem(ctx, ZeroSolveSet, uidS).Err(); err != nil {
				return err
			}
			return nil
		})
		return err
	}
	return doTransactional(ctx, txf, InfoKey)
}
func (acm *ACMStandings) Accept(uid uint, i int, submitTime, startTime time.Time) error {
	defer log.Logger.Sync()
	defer log.LoggerSugar.Sync()

	uidS := strconv.FormatUint(uint64(uid), 10)
	InfoKey := acm.InfoPrefix + uidS

	ctx := context.TODO()
	txf := func(tx *redis.Tx) error {
		result, err := tx.Get(ctx, InfoKey).Result()
		if err != nil {
			log.LoggerSugar.Errorf("user(%v) not register contest(%v)", uid, acm.BindCID)
			return err
		}
		var problemSolveStatus []ProblemSolveStatus
		if err := json.Unmarshal([]byte(result), &problemSolveStatus); err != nil {
			log.Logger.Error("json Unmarshal error", log.Any("err", err))
			return err
		}

		if problemSolveStatus[i].Status == model.AC {
			return nil
		}
		problemSolveStatus[i].Status = model.AC
		d := submitTime.Sub(startTime).Minutes()
		problemSolveStatus[i].Penalty += int64(d)

		// update
		ACCnt, totPenalty := calc(problemSolveStatus)
		oldKey := acm.SolveSetPrefix + strconv.Itoa(ACCnt-1)
		newKey := acm.SolveSetPrefix + strconv.Itoa(ACCnt)

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if err := pipe.ZRem(ctx, oldKey, uidS).Err(); err != nil {
				return err
			}

			if err := pipe.ZAdd(
				ctx,
				newKey,
				redis.Z{Score: float64(totPenalty), Member: uidS},
			).Err(); err != nil {
				return err
			}

			pssJson, err := json.Marshal(problemSolveStatus)
			if err != nil {
				log.Logger.Error("json Marshal error", log.Any("err", err))
				return err
			}
			return pipe.Set(ctx, InfoKey, pssJson, 0).Err()
		})
		return err
	}
	return doTransactional(ctx, txf, InfoKey)
}
func (acm *ACMStandings) Fail(uid uint, i int) error {
	defer log.Logger.Sync()
	defer log.LoggerSugar.Sync()

	uidS := strconv.FormatUint(uint64(uid), 10)
	InfoKey := acm.InfoPrefix + uidS

	ctx := context.TODO()
	txf := func(tx *redis.Tx) error {
		result, err := tx.Get(ctx, InfoKey).Result()
		if err != nil {
			log.LoggerSugar.Errorf("user(%v) not register contest(%v)", uid, acm.BindCID)
			return err
		}
		var problemSolveStatus []ProblemSolveStatus
		if err := json.Unmarshal([]byte(result), &problemSolveStatus); err != nil {
			log.Logger.Error("json Unmarshal error", log.Any("err", err))
			return err
		}

		if problemSolveStatus[i].Status == model.AC {
			return nil
		}
		// update
		problemSolveStatus[i].FailCnt++
		pssJson, err := json.Marshal(problemSolveStatus)
		if err != nil {
			log.Logger.Error("json Marshal error", log.Any("err", err))
			return err
		}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			return pipe.Set(ctx, InfoKey, pssJson, 0).Err()
		})
		return err
	}
	return doTransactional(ctx, txf, InfoKey)
}
func (acm *ACMStandings) GetStandingsByRank(startIdx int64, len int64) ([][]ProblemSolveStatus, error) {
	queryL, queryR := startIdx, startIdx+len

	var tot int64
	var uids []string
	ctx := context.TODO()
	txf := func(tx *redis.Tx) error {
		for i := acm.ProblemsCnt; i >= 0 && tot+1 <= queryR; i-- {
			sz, err := tx.ZCard(ctx, acm.ZKeys[i]).Result()
			if err != nil {
				return err
			}
			L, R := tot+1, tot+sz
			if max(L, queryL) <= min(R, queryR) {
				length := min(R, queryR) - max(L, queryL) + 1
				start := max(L, queryL) - tot - 1
				solveSet := acm.ZKeys[i]
				Zs, err := tx.ZRangeWithScores(ctx, solveSet, start, length).Result()
				if err != nil {
					return err
				}
				for _, z := range Zs {
					uidS, ok := z.Member.(string)
					if !ok {
						return errors.New("assert error")
					}
					uids = append(uids, uidS)
				}
			}
			tot += sz
		}
		return nil
	}
	if err := doTransactional(ctx, txf, acm.ZKeys...); err != nil {
		return nil, err
	}

	ctx = context.TODO()
	cmds, err := redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, uid := range uids {
			pipe.Get(ctx, uid)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	res := make([][]ProblemSolveStatus, 0)
	for _, cmd := range cmds {
		pssJson := cmd.(*redis.StringCmd).Val()
		var problemSolveStatus []ProblemSolveStatus
		if err := json.Unmarshal([]byte(pssJson), &problemSolveStatus); err != nil {
			log.Logger.Error("json Unmarshal error", log.Any("err", err))
			return nil, err
		}
		res = append(res, problemSolveStatus)
	}
	return res, nil
}
func (acm *ACMStandings) Close() error {
	// 1. final standing store(rankList)
	ctx := context.TODO()
	for i := acm.ProblemsCnt; i >= 0; i-- {
		Zs, err := redis.ZRangeWithScores(ctx, acm.ZKeys[i], 0, -1)
		if err != nil {
			return err
		}
		cmds, err := redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			for _, z := range Zs {
				uidS, ok := z.Member.(string)
				if !ok {
					return errors.New("assert error")
				}
				pipe.Get(ctx, uidS)
			}
			return nil
		})
		if err != nil {
			return err
		}

		datas := make([][]ProblemSolveStatus, 0)
		for _, cmd := range cmds {
			pssJson := cmd.(*redis.StringCmd).Val()
			var problemSolveStatus []ProblemSolveStatus
			if err := json.Unmarshal([]byte(pssJson), &problemSolveStatus); err != nil {
				log.Logger.Error("json Unmarshal error", log.Any("err", err))
				return err
			}
			datas = append(datas, problemSolveStatus)
		}

		// mongo update
		filter := bson.M{"_id": acm.BindCID}
		update := bson.M{"$push": bson.M{"rankList": datas}}
		if _, err = dao.GetContestColl().UpdateOne(
			context.Background(),
			filter,
			update,
		); err != nil {
			return err
		}
	}

	// 2. redis key clear
	ctx = context.TODO()
	iter := redis.Scan(ctx, 0, acm.InfoPrefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if err := redis.Del(ctx, key); err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}

	ctx = context.TODO()
	iter = redis.Scan(ctx, 0, acm.SolveSetPrefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if err := redis.Del(ctx, key); err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func doTransactional(ctx context.Context, txf func(*redis.Tx) error, keys ...string) error {
	for i := 0; i < maxRetries; i++ {
		err := redis.Watch(ctx, txf, keys...)
		if err == nil {
			// Success.
			return nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return err
	}

	return errors.New("increment reached maximum number of retries")
}
