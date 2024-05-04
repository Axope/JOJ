package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/common/response"
	"github.com/Axope/JOJ/internal/dao"
	"github.com/Axope/JOJ/internal/manager"
	"github.com/Axope/JOJ/internal/model"
	"github.com/Axope/JOJ/internal/model/contest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type contestService struct {
}

var ContestService = new(contestService)

// func (c *contestService) GetContestList(req *request.GetContestListRequest, uid uint) ([]response.SimpleContest, error) {
// 	findOptions := options.Find().SetLimit(req.Length).SetSkip(req.StartIndex - 1)
// 	findOptions.SetProjection(bson.D{
// 		{Key: "_id", Value: 1},
// 		{Key: "title", Value: 1},
// 		{Key: "status", Value: 1},
// 		{Key: "startTime", Value: 1},
// 		{Key: "duration", Value: 1},
// 		{Key: "note", Value: 1},
// 		{Key: "registered", Value: 1},
// 	})
// 	findOptions.SetSort(map[string]int{"_id": -1})

// 	cursor, err := dao.GetContestColl().Find(context.TODO(), bson.D{{}}, findOptions)
// 	if err != nil {
// 		log.Logger.Debug("find error")
// 		return nil, err
// 	}
// 	defer cursor.Close(context.TODO())

// 	var results []response.SimpleContest
// 	for cursor.Next(context.TODO()) {
// 		var simpleContest response.SimpleContest
// 		if err := cursor.Decode(&simpleContest); err != nil {
// 			log.Logger.Debug("Decode error")
// 			return nil, err
// 		}

// 		// 查询 registered 字段值
// 		var registeredIDs response.RegisteredIDs
// 		if err := cursor.Decode(&registeredIDs); err != nil {
// 			log.Logger.Debug("decode registeredIDs error")
// 			return nil, err
// 		}

// 		var registered bool
// 		for _, id := range registeredIDs.Ids {
// 			if id == uid {
// 				registered = true
// 				break
// 			}
// 		}

// 		// 设置 Registered 字段
// 		simpleContest.IsRegistered = registered

//			results = append(results, simpleContest)
//		}
//		return results, nil
//	}
func (c *contestService) GetContestList(req *request.GetContestListRequest, uid uint) ([]response.SimpleContest, error) {
	findOptions := options.Find().SetLimit(req.Length).SetSkip(req.StartIndex - 1)
	findOptions.SetProjection(bson.D{
		{Key: "_id", Value: 1},
		{Key: "title", Value: 1},
		{Key: "status", Value: 1},
		{Key: "startTime", Value: 1},
		{Key: "duration", Value: 1},
		{Key: "note", Value: 1},
		{Key: "registered", Value: 1},
	})
	findOptions.SetSort(map[string]int{"_id": -1})

	cursor, err := dao.GetContestColl().Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Logger.Debug("find error")
		return nil, err
	} else {
		var results []response.SimpleContest
		if err = cursor.All(context.TODO(), &results); err != nil {
			return nil, err
		}
		log.Logger.Debug("All decode", log.Any("result", results), log.Any("uid", uid))
		for i, simpleContest := range results {
			isRegistered := false
			for _, id := range simpleContest.Registered {
				if id == uid {
					isRegistered = true
					break
				}
			}
			results[i].IsRegistered = isRegistered
		}
		return results, nil
	}
}

func checkStart(c *contest.Contest) error {
	if c.Status == contest.REGISTER {
		return errors.New("contest has not started")
	}
	return nil
}
func (c *contestService) GetContest(cidS string) (*contest.Contest, error) {
	cid, err := primitive.ObjectIDFromHex(cidS)
	if err != nil {
		return nil, err
	}

	findOptions := options.FindOne()
	findOptions.SetProjection(bson.D{{Key: "testCases", Value: 0}})

	var result contest.Contest
	err = dao.GetContestColl().FindOne(context.TODO(), bson.D{{Key: "_id", Value: cid}}, findOptions).Decode(&result)
	if err != nil {
		return nil, err
	} else {
		if err = checkStart(&result); err != nil {
			return nil, err
		}
		return &result, nil
	}
}

func (c *contestService) CreateContest(req *request.CreateContestRequest) (*contest.Contest, error) {
	cid := primitive.NewObjectID()

	var problems []contest.ContestProblem
	if err := json.Unmarshal([]byte(req.ProblemsJson), &problems); err != nil {
		log.Logger.Error("Unmarshal error", log.Any("json", req.ProblemsJson))
		return nil, err
	}
	newContest, err := contest.NewContest(cid, req.Title, problems,
		req.StartTime, time.Minute*time.Duration(req.Duration), req.Note, req.Rule)
	if err != nil {
		return nil, err
	}

	log.Logger.Debug("new contest", log.Any("newContest", newContest))
	// mongo
	insertResult, err := dao.GetContestColl().InsertOne(context.TODO(), newContest)
	if err != nil {
		return nil, err
	}

	manager.ContestManager.NewContest(newContest)

	// TODO: delete
	if insertResult.InsertedID.(primitive.ObjectID) != cid {
		panic("unknown error: cid not equal")
	}
	return newContest, nil
}

func (c *contestService) RegisterContest(req *request.RegisterContestRequest) error {
	cid, err := primitive.ObjectIDFromHex(req.CID)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: cid}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "registered", Value: req.UID},
		}},
	}
	_, err = dao.GetContestColl().UpdateOne(context.TODO(), filter, update)
	return err
}
func (c *contestService) UnregisterContest(req *request.UnregisterContestRequest) error {
	cid, err := primitive.ObjectIDFromHex(req.CID)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: cid}}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "registered", Value: req.UID},
		}},
	}
	_, err = dao.GetContestColl().UpdateOne(context.TODO(), filter, update)
	return err
}
func (c *contestService) GetStandingsByRank(req *request.GetStandingsByRankRequest) ([]contest.RankListData, string, error) {
	defer log.Logger.Sync()
	defer log.LoggerSugar.Sync()
	cid, err := primitive.ObjectIDFromHex(req.Cid)
	if err != nil {
		return nil, "", err
	}
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.D{
		{Key: "rankList", Value: 1},
		{Key: "status", Value: 1},
		{Key: "rule", Value: 1},
	})
	var dbContest contest.Contest
	err = dao.GetContestColl().FindOne(context.TODO(), bson.D{{Key: "_id", Value: cid}}, findOptions).Decode(&dbContest)
	if err != nil {
		log.Logger.Error("FindOne error", log.Any("_id", req.Cid), log.Any("findOptions", findOptions))
		return nil, "", err
	}
	log.LoggerSugar.Debugf("db query contest:%+v", dbContest)

	if dbContest.Status == contest.REGISTER {
		panic("contest is register")
	} else if dbContest.Status == contest.RUNNING {
		managerContest, err := manager.ContestManager.GetContest(req.Cid)
		if err != nil {
			return nil, "", err
		}
		data, err := managerContest.GetStandingsByRank(req.StartIdx, req.Len)
		return data, dbContest.Rule, err
	}
	length := int64(len(dbContest.RankList))
	return dbContest.RankList[req.StartIdx-1 : min(length, req.StartIdx-1+req.Len)], dbContest.Rule, nil
}

func (c *contestService) GetContestSubmissionList(req *request.GetContestSubmissionListRequest) ([]model.Submission, error) {
	cid, err := primitive.ObjectIDFromHex(req.CID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "cid", Value: cid}}
	log.LoggerSugar.Debugf("cid = %v", cid)
	cursor, err := dao.GetSubmissionColl().Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	} else {
		var results []model.Submission
		if err = cursor.All(context.TODO(), &results); err != nil {
			return nil, err
		}
		return results, nil
	}
}
