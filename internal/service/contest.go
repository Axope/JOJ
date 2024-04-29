package service

import (
	"context"
	"errors"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/common/response"
	"github.com/Axope/JOJ/internal/dao"
	"github.com/Axope/JOJ/internal/model/contest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type contestService struct {
}

var ContestService = new(contestService)

func (c *contestService) GetContestList(req *request.GetContestListRequest) ([]response.SimpleContest, error) {
	findOptions := options.Find().SetLimit(req.Length).SetSkip(req.StartIndex - 1)
	findOptions.SetProjection(bson.D{
		{Key: "_id", Value: 1},
		{Key: "title", Value: 1},
		{Key: "status", Value: 1},
		{Key: "startTime", Value: 1},
		{Key: "duration", Value: 1},
		{Key: "note", Value: 1},
	})

	cursor, err := dao.GetContestColl().Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	} else {
		var results []response.SimpleContest
		if err = cursor.All(context.TODO(), &results); err != nil {
			return nil, err
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

	newContest, err := contest.NewContest(cid, req.Title, req.Problems,
		req.StartTime, req.Duration, req.Note, req.Rule)
	if err != nil {
		return nil, err
	}

	log.Logger.Debug("new contest", log.Any("newContest", newContest))
	// mongo
	insertResult, err := dao.GetContestColl().InsertOne(context.TODO(), newContest)
	if err != nil {
		return nil, err
	}

	// TODO: delete
	if insertResult.InsertedID.(primitive.ObjectID) != cid {
		panic("unknown error: cid not equal")
	}

	return newContest, nil
}
