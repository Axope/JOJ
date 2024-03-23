package service

import (
	"context"

	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/internal/dao"
	"github.com/Axope/JOJ/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type problemService struct {
}

var ProblemService = new(problemService)

func (p *problemService) GetProblemList(req *request.GetProblemListRequest) ([]model.Problem, error) {
	findOptions := options.Find().SetLimit(req.Length).SetSkip(req.StartIndex - 1)
	findOptions.SetProjection(bson.D{{Key: "testCases", Value: 0}})

	cursor, err := dao.GetProblemColl().Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	} else {
		var results []model.Problem
		if err = cursor.All(context.TODO(), &results); err != nil {
			return nil, err
		}
		return results, nil
	}
}

func (p *problemService) GetProblem(req *request.GetProblemRequest) (*model.Problem, error) {
	pid, err := primitive.ObjectIDFromHex(req.PID)
	if err != nil {
		return nil, err
	}
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.D{{Key: "testCases", Value: 0}})

	var result model.Problem
	err = dao.GetProblemColl().FindOne(context.TODO(), bson.D{{Key: "_id", Value: pid}}, findOptions).Decode(&result)
	if err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}
